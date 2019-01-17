// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package user

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	randomdata "github.com/Pallinder/go-randomdata"
)

const (
	url = "http://127.0.0.1:8080/v1/users"
)

// MakePostRequest performs an HTTP Post request with JSON payload.
func MakePostRequest(data map[string]interface{}) (*http.Response, error) {
	buff, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buff))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TestPostUsers(t *testing.T) {
	profile := randomdata.GenerateProfile(randomdata.Male | randomdata.Female | randomdata.RandomGender)
	user := map[string]interface{}{
		"username": profile.Login.Username,
		"password": profile.Login.Md5,
		"email":    profile.Email,
	}

	resp, err := MakePostRequest(user)
	if err != nil {
		t.Errorf("TestPostUsers() failed, err: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Errorf("TestPostUsers() failed, got status != 201: %s", resp.Status)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)

}

func TestGetUsers(t *testing.T) {
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("TestGetUsers() failed, err: %s", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("TestPostUsers() failed, got status != 200: %s", resp.Status)
	}

	defer resp.Body.Close()

	var result map[string] interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
}

//

// log.Println(result)
// log.Println(result["data"])
