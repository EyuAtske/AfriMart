package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Make this a default constant, but allow overriding via the struct
const DefaultChapaBaseURL = "https://api.chapa.global/v2"

type Customer struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type Meta struct {
	OrderID string `json:"order_id"`
	Notes   string `json:"notes,omitempty"`
}

type InitializePaymentRequest struct {
	Amount            float64  `json:"amount"`
	Currency          string   `json:"currency"`
	MerchantReference string   `json:"tx_ref"`
	CallbackURL       string   `json:"callback_url"`
	Customer          Customer `json:"customer"`
	Meta              Meta     `json:"meta,omitempty"`
	ReturnURL         string  `json:"return_url"`
}

type InitializePaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference   string `json:"reference"`
		CheckoutURL string `json:"checkout_url"`
		CreatedAt   string `json:"created_at"`
		ExpiresAt   string `json:"expires_at"`
	} `json:"data"`
}

type VerifyPaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference         string  `json:"reference"`
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		Status            string  `json:"status"`
		MerchantReference string  `json:"merchant_reference"`
	} `json:"data"`
}

type ChapaClient struct {
	SecretKey   string
	CallbackURL string
	BaseURL     string 
	HTTPClient  *http.Client
}

func NewClient(secretKey, callbackURL string) *ChapaClient {
	return &ChapaClient{
		SecretKey:   secretKey,
		CallbackURL: callbackURL,
		BaseURL:     DefaultChapaBaseURL, 
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ChapaClient) InitializePayment(ctx context.Context, req InitializePaymentRequest) (*InitializePaymentResponse, error) {
	if c.SecretKey == "" {
		return nil, fmt.Errorf("CHAPA_SECRET_KEY is not configured. Please set it in your environment variables")
	}
	url := c.BaseURL + "/payments/hosted"
	req.CallbackURL = c.CallbackURL 

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.SecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("chapa api error (status %d): %s", resp.StatusCode, errResp.Message)
	}

	var result InitializePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("chapa initialization failed: %s", result.Message)
	}

	return &result, nil
}

func (c *ChapaClient) VerifyPayment(ctx context.Context, reference string) (*VerifyPaymentResponse, error) {
	if c.SecretKey == "" {
		return nil, fmt.Errorf("CHAPA_SECRET_KEY is not configured. Please set it in your environment variables")
	}
	url := fmt.Sprintf("%s/payments/%s/verify", c.BaseURL, reference)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.SecretKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("chapa verify api error (status %d): %s", resp.StatusCode, errResp.Message)
	}

	var result VerifyPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("chapa verification failed: %s", result.Message)
	}

	return &result, nil
}