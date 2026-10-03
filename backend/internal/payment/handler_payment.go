package payment

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/EyuAtske/AfriMart/backend/config"
	"github.com/EyuAtske/AfriMart/backend/internal/comm"
	"github.com/google/uuid"
)

const maxChapaReferenceLen = 20

func newTransactionRef() string {
	const prefix = "AFR-"
	hexID := strings.ReplaceAll(uuid.New().String(), "-", "")
	return prefix + hexID[:maxChapaReferenceLen-len(prefix)]
}

type PaymentHandler struct {
	Config      *config.ApiConfig
	Queries     PaymentQuerier
	Logger      *slog.Logger
	ChapaClient *ChapaClient
}

func NewPaymentHandler(cfg *config.ApiConfig, queries PaymentQuerier, logger *slog.Logger) *PaymentHandler {
	secretKey := os.Getenv("CHAPA_SECRET_KEY")
	callbackURL := os.Getenv("CHAPA_CALLBACK_URL")

	if secretKey == "" {
		logger.Warn("CHAPA_SECRET_KEY is not set in environment")
	}

	return &PaymentHandler{
		Config:      cfg,
		Queries:     queries,
		Logger:      logger,
		ChapaClient: NewClient(secretKey, callbackURL),
	}
}

type CheckoutRequest struct {
	PaymentMethod   string `json:"payment_method"`
	RecipientName   string `json:"recipient_name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	DeliveryAddress string `json:"delivery_address"`
	DeliveryCity    string `json:"delivery_city"`
	DeliveryNotes   string `json:"delivery_notes"`
}

type CheckoutResponse struct {
	PaymentMethod string    `json:"payment_method"`
	OrderID       uuid.UUID `json:"order_id"`
	TransactionID string    `json:"transaction_id"`
	CheckoutURL   string    `json:"checkout_url,omitempty"`
	Message       string    `json:"message"`
}

func writePaymentJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// HandleCallback replaces the old HandleCallback in handler_payment.go.
// PaymentHandler must now have a DB *sql.DB field.
func (h *PaymentHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reference := r.URL.Query().Get("reference")
	if reference == "" {
		reference = r.URL.Query().Get("trx_ref")
	}
	if reference == "" {
		h.Logger.WarnContext(ctx, "callback failed: missing reference")
		comm.RespondErrorWithJson(w, r, http.StatusBadRequest, "Missing payment reference", nil)
		return
	}

	pay, err := h.Queries.GetPaymentByProviderRef(ctx, sql.NullString{String: reference, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.WarnContext(ctx, "callback failed: payment record not found", "reference", reference)
			comm.RespondErrorWithJson(w, r, http.StatusNotFound, "Payment record not found", err)
			return
		}
		h.Logger.ErrorContext(ctx, "callback failed: db error", "reference", reference, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Database error", err)
		return
	}

	// Already paid: nothing more to do.
	if pay.PaymentStatus == "successful" {
		writePaymentJSON(w, http.StatusOK, map[string]string{
			"status": "successful", "message": "Payment already processed", "order_id": pay.OrderID.String(),
		})
		return
	}

	// Always verify with Chapa (even if we already marked it failed/cancelled)
	// so a late successful payment is detected instead of silently ignored.
	verifyResp, err := h.ChapaClient.VerifyPayment(ctx, reference)
	if err != nil {
		h.Logger.ErrorContext(ctx, "chapa verification failed", "reference", reference, "error", err)
		comm.RespondErrorWithJson(w, r, http.StatusBadGateway, "Failed to verify payment with gateway", err)
		return
	}

	if verifyResp.Data.MerchantReference != pay.TransactionID.String {
		h.Logger.ErrorContext(ctx, "merchant reference mismatch",
			"expected", pay.TransactionID.String, "got", verifyResp.Data.MerchantReference)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "Payment verification failed: reference mismatch", nil)
		return
	}

	if got := fmt.Sprintf("%.2f", verifyResp.Data.Amount); got != pay.Amount {
		h.Logger.ErrorContext(ctx, "amount mismatch", "expected", pay.Amount, "got", got)
		comm.RespondErrorWithJson(w, r, http.StatusForbidden, "Payment verification failed: amount mismatch", nil)
		return
	}

	switch verifyResp.Data.Status {
	case "success":
		applied, err := ConfirmPaidOrder(ctx, h.Config.DB, pay.OrderID, pay.ID, reference)
		if err != nil {
			h.Logger.ErrorContext(ctx, "failed to confirm paid order", "payment_id", pay.ID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error updating payment status", err)
			return
		}
		if applied {
			writePaymentJSON(w, http.StatusOK, map[string]string{
				"status": "successful", "message": "Payment successful. Order confirmed.", "order_id": pay.OrderID.String(),
			})
			return
		}

		// Not applied: either a concurrent callback already confirmed it, or the
		// order was cancelled/expired before the customer's payment arrived.
		cur, getErr := h.Queries.GetPaymentByProviderRef(ctx, sql.NullString{String: reference, Valid: true})
		if getErr == nil && cur.PaymentStatus == "successful" {
			writePaymentJSON(w, http.StatusOK, map[string]string{
				"status": "successful", "message": "Payment already processed", "order_id": pay.OrderID.String(),
			})
			return
		}
		h.Logger.ErrorContext(ctx,
			"CRITICAL: Chapa reports success but payment is no longer pending; manual refund/review needed",
			"payment_id", pay.ID, "order_id", pay.OrderID, "reference", reference, "local_status", pay.PaymentStatus)
		writePaymentJSON(w, http.StatusConflict, map[string]string{
			"status":   pay.PaymentStatus,
			"message":  "Payment received but the order is no longer active. Support will contact you.",
			"order_id": pay.OrderID.String(),
		})

	case "failed", "cancelled":
		if pay.PaymentStatus != "pending" {
			writePaymentJSON(w, http.StatusOK, map[string]string{
				"status": pay.PaymentStatus, "message": "Payment already processed", "order_id": pay.OrderID.String(),
			})
			return
		}
		applied, err := CancelUnpaidOrder(ctx, h.Config.DB, pay.OrderID, pay.ID, verifyResp.Data.Status, verifyResp.Message)
		if err != nil {
			h.Logger.ErrorContext(ctx, "failed to cancel unpaid order", "payment_id", pay.ID, "error", err)
			comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Error updating payment status", err)
			return
		}
		msg := "Payment was not successful. Order cancelled."
		if !applied {
			msg = "Payment already processed"
		}
		writePaymentJSON(w, http.StatusPaymentRequired, map[string]string{
			"status": verifyResp.Data.Status, "message": msg, "order_id": pay.OrderID.String(),
		})

	case "pending":
		writePaymentJSON(w, http.StatusOK, map[string]string{
			"status": "pending", "message": "Payment is still being processed",
		})

	default:
		h.Logger.WarnContext(ctx, "unknown chapa status", "status", verifyResp.Data.Status, "reference", reference)
		comm.RespondErrorWithJson(w, r, http.StatusInternalServerError, "Unknown payment status", nil)
	}
}
