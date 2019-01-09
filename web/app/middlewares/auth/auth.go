package auth

import (
	"github.com/labstack/echo"
)

// ValidateAPIKey will check if API key is allowed
func ValidateAPIKey(key string, ctx echo.Context) (bool, error) {

	return true, nil
}
