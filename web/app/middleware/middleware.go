// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

var (
	// RequireLogin check JWT token.
	RequireLogin echo.MiddlewareFunc 
)

// Init middlewares
func Init(e *echo.Echo){

	// logging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	// cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// jwt
	key := viper.GetString("auth.signkey")
	RequireLogin = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(key),
		TokenLookup: "cookie:JWTCookie",
	})

}