// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var (
	// RequireLogin check JWT token.
	RequireLogin echo.MiddlewareFunc

	// RequireToken checks email confirmation token.
	RequireToken echo.MiddlewareFunc
)

// CustomClaims are custom claims extending default ones.
// Used for email confirmation and pwd reset.
type CustomClaims struct {
	Username string `json:"username"`
	Purpose  string `json:"purpose"`
	jwt.StandardClaims
}

// LoginCustomClaims are custom claims extending default ones
type LoginCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))

	// cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{viper.GetString("ui.address")},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))

	// trailing slash
	e.Pre(middleware.AddTrailingSlash())

	// jwt
	key := viper.GetString("auth.signkey")
	RequireLogin = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(key),
		TokenLookup: "cookie:JWTCookie",
	})

	RequireToken = middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &CustomClaims{},
		SigningKey:  []byte(key),
		TokenLookup: "query:token",
	})
}
