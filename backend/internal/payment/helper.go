package payment

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type PendingOrderExpiryJob struct {
	DB       *sql.DB
	Queries  *database.Queries
	Logger   *slog.Logger
	Interval time.Duration
	Timeout  time.Duration
}

func CancelUnpaidOrder(
	ctx context.Context,
	db *sql.DB,
	orderID, paymentID uuid.UUID,
	newStatus, reason string,
) (applied bool, err error) {
	if newStatus != "failed" && newStatus != "cancelled" {
		return false, fmt.Errorf("invalid payment status %q for cancellation", newStatus)
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin compensation tx: %w", err)
	}
	defer tx.Rollback()

	q := database.New(tx)

	rows, err := q.UpdatePaymentStatusConditional(ctx, database.UpdatePaymentStatusConditionalParams{
		ID:            paymentID,
		CurrentStatus: "pending",
		NewStatus:     newStatus,
		FailureReason: sql.NullString{String: reason, Valid: reason != ""},
	})
	if err != nil {
		return false, fmt.Errorf("update payment status: %w", err)
	}
	if rows == 0 {
		return false, nil
	}

	if _, err := q.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     orderID,
		Status: "cancelled",
	}); err != nil {
		return false, fmt.Errorf("cancel order: %w", err)
	}

	items, err := q.GetOrderItems(ctx, orderID)
	if err != nil {
		return false, fmt.Errorf("load order items: %w", err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ProductID.String() < items[j].ProductID.String()
	})

	for _, item := range items {
		if _, err := q.RestoreProductStock(ctx, database.RestoreProductStockParams{
			ID:    item.ProductID,
			Stock: item.Quantity,
		}); err != nil {
			return false, fmt.Errorf("restore stock for product %s: %w", item.ProductID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit compensation tx: %w", err)
	}
	return true, nil
}

func ConfirmPaidOrder(
	ctx context.Context,
	db *sql.DB,
	orderID, paymentID uuid.UUID,
	providerRef string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin confirm tx: %w", err)
	}
	defer tx.Rollback()

	q := database.New(tx)

	rows, err := q.UpdatePaymentStatusConditional(ctx, database.UpdatePaymentStatusConditionalParams{
		ID:                paymentID,
		CurrentStatus:     "pending",
		NewStatus:         "successful",
		PaidAt:            sql.NullTime{Time: time.Now(), Valid: true},
		ProviderReference: sql.NullString{String: providerRef, Valid: providerRef != ""},
	})
	if err != nil {
		return false, fmt.Errorf("update payment status: %w", err)
	}
	if rows == 0 {
		return false, nil
	}

	if _, err := q.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     orderID,
		Status: "confirmed",
	}); err != nil {
		return false, fmt.Errorf("confirm order: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit confirm tx: %w", err)
	}
	return true, nil
}

func NewPendingOrderExpiryJob(db *sql.DB, queries *database.Queries, logger *slog.Logger) *PendingOrderExpiryJob {
	return &PendingOrderExpiryJob{
		DB:       db,
		Queries:  queries,
		Logger:   logger,
		Interval: 5 * time.Minute,  // Check every 5 minutes
		Timeout:  30 * time.Minute, // Cancel if pending for > 30 mins
	}
}

func (j *PendingOrderExpiryJob) Start(ctx context.Context) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()

	j.Logger.Info("pending order expiry job started", "interval", j.Interval, "timeout", j.Timeout)
	j.run(ctx) // Run once immediately

	for {
		select {
		case <-ctx.Done():
			j.Logger.Info("pending order expiry job stopping")
			return
		case <-ticker.C:
			j.run(ctx)
		}
	}
}

func (j *PendingOrderExpiryJob) run(ctx context.Context) {
	cutoff := time.Now().Add(-j.Timeout)

	expired, err := j.Queries.GetExpiredPendingPayments(ctx, cutoff)
	if err != nil {
		j.Logger.ErrorContext(ctx, "expiry job: failed to fetch expired payments", "error", err)
		return
	}

	if len(expired) == 0 {
		return
	}

	j.Logger.InfoContext(ctx, "expiry job: found expired pending payments", "count", len(expired))

	for _, p := range expired {
		// Use a detached context with a timeout so the job doesn't hang
		cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		applied, err := CancelUnpaidOrder(cancelCtx, j.DB, p.OrderID, p.ID, "cancelled", "Payment checkout expired")
		cancel()

		if err != nil {
			j.Logger.ErrorContext(ctx, "expiry job: failed to cancel expired order",
				"order_id", p.OrderID, "payment_id", p.ID, "error", err)
			continue
		}

		if applied {
			j.Logger.InfoContext(ctx, "expiry job: successfully cancelled expired order",
				"order_id", p.OrderID, "payment_id", p.ID)
		}
	}
}
