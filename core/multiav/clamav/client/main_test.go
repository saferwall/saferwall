// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package clamav

import (
	"log"
	"testing"
)

func TestScanFile(t *testing.T) {
	clamclient := Init()
	clamres, _ := ScanFile(clamclient, "eicar")
	log.Println(clamres)

}
