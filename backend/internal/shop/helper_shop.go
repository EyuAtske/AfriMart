package shop

import (
	"errors"
	"strings"
)

func validateCreateShopRequest(params *createShopRequest) error {
	params.Name = strings.TrimSpace(params.Name)

	if params.Name == "" {
		return errors.New("shop name is required")
	}

	return nil
}
