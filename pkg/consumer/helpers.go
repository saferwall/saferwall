// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"errors"
)

func loadConfig() {
	viper.AddConfigPath("configs")

	// Load the config type depending on env variable.
	var name string
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		name = "saferwall.dev"
	case "prod":
		name = "saferwall.prod"
	case "test":
		name = "saferwall.test"
	default:
		log.Fatal("ENVIRONMENT is not set")
	}

	viper.SetConfigName(name) // no need to include file extension
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Config %s was loaded", name)
}

func login() string {
	username := viper.GetString("backend.admin_user")
	password := viper.GetString("backend.admin_pwd")
	requestBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	url := viper.GetString("backend.address") + "/v1/auth/login"
	body := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	var res map[string]string
	json.Unmarshal(d, &res)
	return res["token"]	
}

func updateDocument(sha256 string, buff []byte) error {
	client := &http.Client{}
	client.Timeout = time.Second * 15
	url := backendEndpoint + sha256
	log.Infoln("Sending results to ", url)

	body := bytes.NewBuffer(buff)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Errorf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Cookie", "JWTCookie="+backendToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("client.Do() failed with '%s'\n", err)
	}

	// check if token is not expired
	if resp.StatusCode == http.StatusUnauthorized {
		backendToken = login()
		req.Header.Set("Cookie", "JWTCookie="+backendToken)
		resp, err = client.Do(req)
		if err != nil {
			log.Errorf("retry: client.Do() failed with '%s'", err)
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return errors.New("Failed to get a new login token")
		}
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	log.Infof("Response status code: %d, text: %s", resp.StatusCode, string(d))
	return nil
}
