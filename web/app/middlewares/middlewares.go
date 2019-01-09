package middlewares

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/saferwall/saferwall/web/server/app/middlewares/auth"
)

func requireJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("content-type")
		fmt.Println(c.Request().Header)
		if contentType != "application/json" {
			return echo.NewHTTPError(http.StatusBadRequest)
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

	// require JSON
	// e.Use(requireJSON)

	// authorization
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key",
		Validator: middleware.KeyAuthValidator(auth.ValidateAPIKey),
	}))

}
