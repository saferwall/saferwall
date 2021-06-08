// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// HTTP Constants
	authPath      = "/v1/auth/login/"
	reqTimeout    = time.Duration(5 * time.Second)
	filePath      = "/v1/files/"
	updateTimeout = time.Duration(15 * time.Second)
)

var (
	errHTTPStatusCodeNotOK    = errors.New("http response status code != 200")
	errHTTPStatusUnauthorized = errors.New("jwt token expired")
)

// fetchAuthToken retrieves a JWT token from the API.
func fetchAuthToken(cfg *Config) (string, error) {

	var authToken string

	authRequest, err := newAuthReq(cfg.Backend.Address, cfg.Backend.Username, cfg.Backend.Password)
	if err != nil {
		return "", err
	}
	httpClient := http.Client{
		Timeout: reqTimeout,
	}
	resp, err := httpClient.Do(authRequest)
	if err != nil {
		return authToken, err
	}

	if resp.StatusCode != http.StatusOK {
		return authToken, errHTTPStatusCodeNotOK
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return authToken, err
	}

	var res map[string]string
	err = json.Unmarshal(d, &res)
	if err != nil {
		return authToken, err
	}

	authToken = res["token"]
	return authToken, nil
}

// newAuthReq builds a new HTTP request for a JWT token.
func newAuthReq(address, username, password string) (*http.Request, error) {
	requestBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	url := address + authPath
	body := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	// Set HTTP Header
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	return request, nil
}

// getAuthToken() retrieves a JWT auth token from the web apis.
func getAuthToken(cfg *Config) (string, error) {

	var authToken string

	requestBody, err := json.Marshal(map[string]string{
		"username": cfg.Backend.Username,
		"password": cfg.Backend.Password,
	})
	if err != nil {
		return authToken, err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	url := cfg.Backend.Address + "/v1/auth/login/"
	body := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return authToken, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(request)
	if err != nil {
		return authToken, err
	}

	if resp.StatusCode != http.StatusOK {
		return authToken, errHTTPStatusCodeNotOK
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return authToken, err
	}

	var res map[string]string
	err = json.Unmarshal(d, &res)
	if err != nil {
		return authToken, err
	}

	authToken = res["token"]
	return authToken, nil
}

func updateDocument(sha256, token string, cfg *Config, buff []byte) error {

	url := cfg.Backend.Address + filePath + sha256
	body := bytes.NewBuffer(buff)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Cookie", "JWTCookie="+token)

	client := http.Client{
		Timeout: updateTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check if token is not expired.
	if resp.StatusCode == http.StatusUnauthorized {
		return errHTTPStatusUnauthorized
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
