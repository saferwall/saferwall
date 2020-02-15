// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package route

import (
	"github.com/labstack/echo/v4"
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

	// handle user login and confirmation via email.
	e.POST("/v1/auth/login/", auth.Login, m.RequireJSON) 	
	e.GET("/v1/auth/confirm/", auth.Confirm, m.RequireToken)
	e.POST("/v1/auth/resend-confirmation/", auth.ReconfirmAccount)

	// To reset the current password (in case user forget the password).
	e.DELETE("/v1/users/password/", auth.ResetPassword, m.RequireJSON)

	// To create new password (if user has reset the password).
	// new password, activation code and emailID should be given in body.
	e.POST("/v1/users/password/", auth.NewPassword, m.RequireToken, m.RequireJSON)

	// To update the password (if user knows is old password and new password)
	e.PUT("/v1/users/:username/password/", auth.UpdatePassword, m.RequireLogin, m.RequireJSON)

	// To update the email (if user knows its password)
	e.PUT("/v1/users/:username/email/", auth.UpdateEmail, m.RequireLogin, m.RequireJSON)

	// To update user avatar
	e.PUT("/v1/users/:username/avatar/", user.UpdateAvatar, m.RequireLogin)

	// handle /files endpoint.
	e.GET("/v1/files/", file.GetFiles, m.RequireLogin, auth.IsAdmin)
	e.POST("/v1/files/", file.PostFiles, m.RequireLogin)
	e.PUT("/v1/files/", file.PutFiles, m.RequireLogin, auth.IsAdmin)
	e.DELETE("/v1/files/", file.DeleteFiles, m.RequireLogin, auth.IsAdmin)

	// handle /files/:sha256 endpoint.
	e.GET("/v1/files/:sha256/", file.GetFile)
	e.PUT("/v1/files/:sha256/", file.PutFile, m.RequireLogin, m.RequireJSON)
	e.DELETE("/v1/files/:sha256/", file.DeleteFile, m.RequireLogin, auth.IsAdmin)

	// handle comments
	e.POST("/v1/files/:sha256/comments/", file.PostComment, m.RequireLogin, m.RequireJSON)
	e.DELETE("/v1/files/:sha256/comments/:id/", file.DeleteComment, m.RequireLogin)

	// handle file download.
	e.POST("/v1/files/:sha256/actions/", file.Actions, m.RequireLogin)
	e.GET("/v1/files/:sha256/download/", file.Download, m.RequireLogin)

	// handle /users endpoint.
	e.GET("/v1/users/", user.GetUsers, m.RequireLogin, auth.IsAdmin)
	e.POST("/v1/users/", user.PostUsers, m.RequireJSON)
	e.PUT("/v1/users/", user.PutUsers, m.RequireLogin, auth.IsAdmin)
	e.DELETE("/v1/users/", user.DeleteUsers, m.RequireLogin, auth.IsAdmin)

	// handle /users/:username endpoint.
	e.GET("/v1/users/:username/", user.GetUser)
	e.PUT("/v1/users/:username/", user.PutUser, m.RequireLogin)
	e.DELETE("/v1/users/:username/", user.DeleteUser, m.RequireLogin)
	e.GET("/v1/users/:username/avatar/", user.GetAvatar)

	// actions over a user
	e.POST("/v1/users/:username/actions/", user.Actions, m.RequireLogin, m.RequireJSON)

	// handle /admin endpoint
	e.GET("/admin/", auth.Admin, m.RequireLogin, auth.IsAdmin)

	// ugly hack
	user.CreateAdminUser()
	
	return e
}
