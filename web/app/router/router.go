package router

import (
	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/server/app/handlers"
	"github.com/saferwall/saferwall/web/server/app/middlewares"
)

// New create an echo insance
func New() *echo.Echo {

	// Create `echo` instance
	e := echo.New()

	// Setup middlwares
	middlewares.Init(e)

	// `/files`` endpoint
	e.GET("/v1/files", handlers.GetFiles)
	e.POST("/v1/files", handlers.PostFiles)
	e.PUT("/v1/files", handlers.PutFiles)
	e.DELETE("/v1/files", handlers.DeleteFiles)

	// `/files/:sha256`
	e.GET("/v1/files/:sha256", handlers.GetFile)
	e.PUT("/v1/files/:sha256", handlers.PutFile)
	e.DELETE("/v1/files/:sha256", handlers.DeleteFile)

	// `/users` endpoint
	e.GET("/v1/users", handlers.GetUsers)
	e.POST("/v1/users", handlers.PostUsers)
	e.PUT("/v1/users", handlers.PutUsers)
	e.DELETE("/v1/users", handlers.DeleteUsers)

	// `/users/:username`
	e.GET("/v1/users/:username", handlers.GetUser)
	e.POST("/v1/users/:username", handlers.PostUser)
	e.PUT("/v1/users/:username", handlers.PutUser)
	e.DELETE("/v1/users/:username", handlers.DeleteUser)

	return e
}
