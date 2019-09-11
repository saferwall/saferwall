// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/route"
	"github.com/spf13/viper"
)

func main() {

	// Init our app
	app.Init()

	// Create echo instance and load all routes
	e := route.New()

	// Start the server
	address := viper.GetString("backend.address")
	e.Logger.Fatal(e.Start(address))
}
