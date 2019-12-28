// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/route"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"

)



func main() {

	// Init our app
	app.Init()

	// Create echo instance and load all routes
	e := route.New()

	address := viper.GetString("app.address")

	// Start the server
	log.Info("Running in debug mode")
	e.Logger.Fatal(e.Start(address))

}
