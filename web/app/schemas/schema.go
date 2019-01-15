// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package schemas

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/xeipuuv/gojsonschema"
)

var (
	// UserSchemaLoader represent a user
	UserSchemaLoader gojsonschema.Schema
)

// Load will load schemas at server startup
func Load() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out := path.Join(dir, "app", "schemas", "user.json")
	out = fmt.Sprintf("file:///%s", out)
	schemaLoader := gojsonschema.NewReferenceLoader(out)
	userSchema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		fmt.Println("Error while loading user schema : ", err)
	}
	UserSchemaLoader = *userSchema
}
