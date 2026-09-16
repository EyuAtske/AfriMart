package shop

import "testing"

func TestValidateCreateShopRequest(t *testing.T) {
	tests := []struct {
		name          string
		requestName   string
		wantName      string
		wantError     bool
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