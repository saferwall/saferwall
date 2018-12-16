package main

import (
	"github.com/saferwall/saferwall/web/server/app"
	"github.com/saferwall/saferwall/web/server/app/common/database"
	"github.com/saferwall/saferwall/web/server/app/router"
	"github.com/saferwall/saferwall/web/server/app/schemas"
	"github.com/saferwall/saferwall/web/server/config"
)

func main() {

	// Load the configuration file
	config.Load()

	// Init our app
	app.Init()

	// Connect to the database
	database.Connect()

	// Load schemas
	schemas.Load()

	// Create echo instance and load all routes
	e := router.New()

	// Start the server
	e.Logger.Fatal(e.Start(":80"))
}
