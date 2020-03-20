// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	uploadURL = "https://api.saferwall.com/v1/files/"
	authURL   = "https://api.saferwall.com/v1/auth/login/"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func newfileUploadRequest(uri, fieldname, filename string) (*http.Request, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func login(username, password string) (string, error) {
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
	body := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(http.MethodPost, authURL, body)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	var res map[string]string
	json.Unmarshal(d, &res)

	if resp.StatusCode != http.StatusOK {
		fmt.Println(res)
		return "", errors.New("Failed login")
	}

	return res["token"], nil
}

func upload(filepath string, authToken string) {

	// Create a new file upload request.
	request, err := newfileUploadRequest(uploadURL, "file", filepath)
	check(err)

	// Add our auth token.
	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	check(err)

	// Read the response.
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
	fmt.Println(body)

}
func main() {

	if len(os.Args) != 2 {
		log.Println("Usage: uploader <filepath>")
		return
	}
	filepath := os.Args[1]
	_, err := os.Stat(filepath)
    if os.IsNotExist(err) {
		log.Fatalf("%s does not exist", filepath)
    }

	// Get credentials.
	username := os.Getenv("SAFERWALL_AUTH_USERNAME")
	password := os.Getenv("SAFERWALL_AUTH_PASSWORD")
	if username == "" || password == "" {
		log.Fatal("SAFERWALL_AUTH_USERNAME or SAFERWALL_AUTH_USERNAME env variable are not set.")
	}

	// Obtain a token.
	token, err := login(username, password)
	check(err)

	// Walk over directory.
	fileList := []string{}
	filepath.Walk(os.Args[1], func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	// Upload files
	for _, file := range fileList {
		upload(file, token)
	}
}
