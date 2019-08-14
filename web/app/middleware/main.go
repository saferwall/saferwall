// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"strings"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

var (
	// RequireLogin check JWT token.
	RequireLogin echo.MiddlewareFunc

	// RequireEmailConfirmationToken checks email confirmation token.
	RequireEmailConfirmationToken echo.MiddlewareFunc
)

// RequireJSON requires an application/json content type.
func RequireJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("content-type")
		if contentType != "" && !strings.Contains(contentType, "application/json") {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"verbose_msg": "Request requires content type: application/json"})
		}
		return next(c)
	}
}

// Init middlewares
func Init(e *echo.Echo) {

	// logging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency_human}` + "\n",
	  }))

	// cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// trailing slash
	e.Pre(middleware.AddTrailingSlash())

	// jwt
	key := viper.GetString("auth.signkey")
	RequireLogin = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(key),
		TokenLookup: "cookie:JWTCookie",
	})

	RequireEmailConfirmationToken = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(key),
		TokenLookup: "query:token",
	})

}
