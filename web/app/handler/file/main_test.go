// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package file

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	url = "http://127.0.0.1:8080/v1/files"
)

func TestPostFiles(t *testing.T) {
	f, err := os.Open("../../../../test/multiav/infected/zbot")
	if err != nil {
		t.Fatalf("TestPostFiles() failed, err: %s", err)
	}
	defer f.Close()

	resp, err := http.Post(url, "binary/octet-stream", f)
	if err != nil {
		t.Errorf("TestPostFiles() failed, err: %s", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
}

func TestGetFiles(t *testing.T) {
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("TestGetUsers() failed, err: %s", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("TestPostUsers() failed, got status != 200: %s", resp.Status)
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
}
