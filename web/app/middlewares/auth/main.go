// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/labstack/echo"
)

// ValidateAPIKey will check if API key is allowed
func ValidateAPIKey(key string, ctx echo.Context) (bool, error) {

	return true, nil
}
