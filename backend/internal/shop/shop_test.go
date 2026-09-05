package shop

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EyuAtske/AfriMart/backend/internal/auth"
	"github.com/EyuAtske/AfriMart/backend/internal/database"
	"github.com/google/uuid"
)

type mockShopQueries struct {
	deactivateShopFunc func(
		ctx context.Context,
		arg database.DeactivateShopParams,
	) (database.Shop, error)

	activateShopFunc func(
		ctx context.Context,
		arg database.ActivateShopParams,
	) (database.Shop, error)

	getShopByIDAndOwnerID func(
    	ctx context.Context,
    	arg database.GetShopByIDAndOwnerIDParams,
	) (database.Shop, error)
}

func (m *mockShopQueries) DeactivateShop(
	ctx context.Context,
	arg database.DeactivateShopParams,
) (database.Shop, error) {
	return m.deactivateShopFunc(ctx, arg)
}

func (m *mockShopQueries) ActivateShop(
	ctx context.Context,
	arg database.ActivateShopParams,
) (database.Shop, error) {
	return m.activateShopFunc(ctx, arg)
}

func (m *mockShopQueries) GetShopByIDAndOwnerID(
	ctx context.Context,
	arg database.GetShopByIDAndOwnerIDParams,
) (database.Shop, error) {
	return m.getShopByIDAndOwnerID(ctx, arg)
}

func TestValidateCreateShopRequest(t *testing.T) {
	tests := []struct {
		name        string
		requestName string
		wantName    string
		wantError   bool
	}{
		{
			name:        "valid shop name",
			requestName: "My Shop",
			wantName:    "My Shop",
			wantError:   false,
		},
		{
			name:        "name with surrounding spaces",
			requestName: "  My Shop  ",
			wantName:    "My Shop",
			wantError:   false,
		},
		{
			name:        "empty name",
			requestName: "",
			wantName:    "",
			wantError:   true,
		},
		{
			name:        "whitespace only name",
			requestName: "     ",
			wantName:    "",
			wantError:   true,
		},
		{
			name:        "tabs and spaces only",
			requestName: " \t  ",
			wantName:    "",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := createShopRequest{
				Name: tt.requestName,
			}

			err := validateCreateShopRequest(&params)

			if tt.wantError && err == nil {
				t.Fatal("Expected validation error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Fatalf(
					"Expected no validation error, got: %v",
					err,
				)
			}

			if params.Name != tt.wantName {
				t.Fatalf(
					"Expected name %q, got %q",
					tt.wantName,
					params.Name,
				)
			}
		})
	}
}

func TestHandleDeactivateShop_Success(t *testing.T) {
	userID := uuid.New()
	shopID := uuid.New()

	expectedShop := database.Shop{
		ID:      shopID,
		OwnerID: userID,
		Name:    "Test Shop",
		Status:  "inactive",
	}

	mockQueries := &mockShopQueries{
		deactivateShopFunc: func(
			ctx context.Context,
			arg database.DeactivateShopParams,
		) (database.Shop, error) {

			if arg.ID != shopID {
				t.Errorf("expected shop ID %v, got %v", shopID, arg.ID)
			}

			if arg.OwnerID != userID {
				t.Errorf("expected owner ID %v, got %v", userID, arg.OwnerID)
			}

			return expectedShop, nil
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/"+shopID.String()+"/deactivate",
		nil,
	)

	req.SetPathValue("shopID", shopID.String())

	ctx := auth.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.HandleDeactivateShop(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var shop database.Shop

	if err := json.NewDecoder(rec.Body).Decode(&shop); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if shop.ID != expectedShop.ID {
		t.Errorf(
			"expected shop ID %v, got %v",
			expectedShop.ID,
			shop.ID,
		)
	}

	if shop.Status != "inactive" {
		t.Errorf(
			"expected status inactive, got %v",
			shop.Status,
		)
	}
}

func TestHandleDeactivateShop_Unauthorized(t *testing.T) {
	shopID := uuid.New()

	mockQueries := &mockShopQueries{
		deactivateShopFunc: func(
			ctx context.Context,
			arg database.DeactivateShopParams,
		) (database.Shop, error) {
			t.Fatal("DeactivateShop should not be called")
			return database.Shop{}, nil
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/"+shopID.String()+"/deactivate",
		nil,
	)

	req.SetPathValue("shopID", shopID.String())

	rec := httptest.NewRecorder()

	handler.HandleDeactivateShop(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestHandleDeactivateShop_InvalidShopID(t *testing.T) {
	userID := uuid.New()

	mockQueries := &mockShopQueries{
		deactivateShopFunc: func(
			ctx context.Context,
			arg database.DeactivateShopParams,
		) (database.Shop, error) {
			t.Fatal("DeactivateShop should not be called")
			return database.Shop{}, nil
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/not-a-uuid/deactivate",
		nil,
	)

	req.SetPathValue("shopID", "not-a-uuid")

	ctx := auth.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.HandleDeactivateShop(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandleDeactivateShop_NotOwner(t *testing.T) {
	userID := uuid.New()
	shopID := uuid.New()

	mockQueries := &mockShopQueries{
		deactivateShopFunc: func(
			ctx context.Context,
			arg database.DeactivateShopParams,
		) (database.Shop, error) {

			if arg.OwnerID != userID {
				t.Errorf(
					"expected owner ID %v, got %v",
					userID,
					arg.OwnerID,
				)
			}

			return database.Shop{}, sql.ErrNoRows
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/"+shopID.String()+"/deactivate",
		nil,
	)

	req.SetPathValue("shopID", shopID.String())

	ctx := auth.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.HandleDeactivateShop(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestHandleActivateShop_Success(t *testing.T) {
	userID := uuid.New()
	shopID := uuid.New()

	expectedShop := database.Shop{
		ID:      shopID,
		OwnerID: userID,
		Name:    "Test Shop",
		Status:  "active",
	}

	mockQueries := &mockShopQueries{
		activateShopFunc: func(
			ctx context.Context,
			arg database.ActivateShopParams,
		) (database.Shop, error) {

			if arg.ID != shopID {
				t.Errorf("expected shop ID %v, got %v", shopID, arg.ID)
			}

			if arg.OwnerID != userID {
				t.Errorf("expected owner ID %v, got %v", userID, arg.OwnerID)
			}

			return expectedShop, nil
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/"+shopID.String()+"/activate",
		nil,
	)

	req.SetPathValue("shopID", shopID.String())

	ctx := auth.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.HandleActivateShop(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var shop database.Shop

	if err := json.NewDecoder(rec.Body).Decode(&shop); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if shop.ID != expectedShop.ID {
		t.Errorf(
			"expected shop ID %v, got %v",
			expectedShop.ID,
			shop.ID,
		)
	}

	if shop.Status != "active" {
		t.Errorf(
			"expected status active, got %v",
			shop.Status,
		)
	}
}

func TestHandleActivateShop_Unauthorized(t *testing.T) {
	shopID := uuid.New()

	mockQueries := &mockShopQueries{
		activateShopFunc: func(
			ctx context.Context,
			arg database.ActivateShopParams,
		) (database.Shop, error){
			t.Fatal("ActivateShop should not be called")
			return database.Shop{}, nil
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/"+shopID.String()+"/activate",
		nil,
	)

	req.SetPathValue("shopID", shopID.String())

	rec := httptest.NewRecorder()

	handler.HandleActivateShop(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestHandleActivateShop_InvalidShopID(t *testing.T) {
	userID := uuid.New()

	mockQueries := &mockShopQueries{
		activateShopFunc: func(
			ctx context.Context,
			arg database.ActivateShopParams,
		) (database.Shop, error) {
			t.Fatal("DeactivateShop should not be called")
			return database.Shop{}, nil
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/not-a-uuid/activate",
		nil,
	)

	req.SetPathValue("shopID", "not-a-uuid")

	ctx := auth.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.HandleActivateShop(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandleActivateShop_NotOwner(t *testing.T) {
	userID := uuid.New()
	shopID := uuid.New()

	mockQueries := &mockShopQueries{
		activateShopFunc: func(
			ctx context.Context,
			arg database.ActivateShopParams,
		) (database.Shop, error) {

			if arg.OwnerID != userID {
				t.Errorf(
					"expected owner ID %v, got %v",
					userID,
					arg.OwnerID,
				)
			}

			return database.Shop{}, sql.ErrNoRows
		},
	}

	handler := &ShopHandler{
		Queries: mockQueries,
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/shops/"+shopID.String()+"/activate",
		nil,
	)

	req.SetPathValue("shopID", shopID.String())

	ctx := auth.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.HandleActivateShop(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}