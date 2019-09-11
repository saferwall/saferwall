// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package file

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	url = "http://localhost:8080/v1/files"
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	log.Println(bodyString)
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

func TestGetFile(t *testing.T) {
	resp, err := http.Get(url + "/df50dd428c2c0a6c2bffc6720b10d690061f1e3e0d1f5ef2f926942cbf4fc69c")
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

