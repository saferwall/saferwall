// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadConfig() {

	// Add config path directories.
	viper.AddConfigPath("configs")
	viper.AddConfigPath("../../configs")

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

func setupLogging() {

	level := viper.GetString("consumer.log_level")
	if len(level) > 0 {
		switch level {
		case "panic":
			log.SetLevel(log.PanicLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "trace":
			log.SetLevel(log.TraceLevel)
		default:
			log.SetLevel(log.WarnLevel)
		}
	} else {
		log.SetLevel(log.WarnLevel)
	}

}

func login() (string, error) {
	username := viper.GetString("backend.admin_user")
	password := viper.GetString("backend.admin_pwd")
	requestBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return "", err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	url := viper.GetString("backend.address") + "/v1/auth/login"
	body := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Errorf("http.NewRequest() failed with: %v", err)
		return "", err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(request)
	if err != nil {
		log.Errorf("client.Do() failed with: %v", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("login() http response status code not 200")
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll() failed with: %v", err)
		return "", err
	}

	var res map[string]string
	err = json.Unmarshal(d, &res)
	if err != nil {
		log.Errorf("json unmarshall failed with: %v", err)
		return "", err
	}
	return res["token"], nil
}

func updateDocument(sha256 string, buff []byte) error {
	client := &http.Client{}
	client.Timeout = time.Second * 15
	url := backendEndpoint + sha256

	body := bytes.NewBuffer(buff)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		contextLogger.Errorf("http.NewRequest() failed with: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Cookie", "JWTCookie="+backendToken)
	resp, err := client.Do(req)
	if err != nil {
		contextLogger.Errorf("client.Do() failed with: %v", err)
		return err
	}

	// Check if token is not expired.
	if resp.StatusCode == http.StatusUnauthorized {
		backendToken, err = login()
		if err != nil {
			contextLogger.Errorf("Failed to get new auth token: %v", err)
			return err
		}
		req.Header.Set("Cookie", "JWTCookie="+backendToken)
		resp, err = client.Do(req)
		if err != nil {
			contextLogger.Errorf("retry: client.Do() failed with: %v", err)
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return errors.New("Failed to get a new login token")
		}
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		contextLogger.Errorf("ioutil.ReadAll() failed with: %v", err)
		return err
	}

	contextLogger.Infof("Scanning finished: status code: %d, resp: %s",
		resp.StatusCode, string(d))
	return nil
}
