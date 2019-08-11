// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package route

import (
	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/app/handler/auth"
	"github.com/saferwall/saferwall/web/app/handler/file"
	"github.com/saferwall/saferwall/web/app/handler/user"
	m "github.com/saferwall/saferwall/web/app/middleware"
)

// New create an echo insance
func New() *echo.Echo {

	// Create `echo` instance
	e := echo.New()

	// Setup middlwares
	m.Init(e)

	// handle /login
	e.POST("/auth/login", auth.Login, m.RequireJSON)
	e.GET("/auth/confirm", auth.Confirm, m.RequireEmailConfirmationToken)

	// handle /files endpoint.
	e.GET("/v1/files", file.GetFiles)
	e.POST("/v1/files", file.PostFiles, m.RequireLogin)
	e.PUT("/v1/files", file.PutFiles, m.RequireLogin)
	e.DELETE("/v1/files", file.DeleteFiles)

	// handle /files/:sha256 endpoint.
	e.GET("/v1/files/:sha256", file.GetFile)
	e.PUT("/v1/files/:sha256", file.PutFile, m.RequireJSON)
	e.DELETE("/v1/files/:sha256", file.DeleteFile)

	// handle /users endpoint.
	e.GET("/v1/users", user.GetUsers)
	e.POST("/v1/users", user.PostUsers, m.RequireJSON)
	e.PUT("/v1/users", user.PutUsers, m.RequireLogin)
	e.DELETE("/v1/users", user.DeleteUsers)

	// handle /users/:username  endpoint.
	e.GET("/v1/users/:username", user.GetUser)
	e.PUT("/v1/users/:username", user.PutUser, m.RequireLogin)
	e.DELETE("/v1/users/:username", user.DeleteUser)

	// handle /admin endpoint
	e.GET("/admin", auth.Admin, m.RequireLogin, auth.IsAdmin)

	return e
}
