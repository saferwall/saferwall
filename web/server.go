// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/common/db"
	"github.com/saferwall/saferwall/web/app/route"
	"github.com/saferwall/saferwall/web/app/schema"
	"github.com/saferwall/saferwall/web/config"
)

func main() {

	// Load the configuration file
	config.Load()

	// Init our app
	app.Init()

	// Connect to the database
	db.Connect()

	// Load schemas
	schema.Load()

	// Create echo instance and load all routes
	e := route.New()

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
