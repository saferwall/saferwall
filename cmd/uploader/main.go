// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"encoding/hex"
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
	"crypto/sha256"
)

const (
	fileURL = "https://api.saferwall.com/v1/files/"
	authURL   = "https://api.saferwall.com/v1/auth/login/"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getSha256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
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
	request, err := newfileUploadRequest(fileURL, "file", filepath)
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

func rescan(sha256 string, authToken string) error {

	payload, err := json.Marshal(map[string]string{
		"type": "rescan",
	})
	if err != nil {
		return err
	}

	url := fileURL + sha256 + "/actions"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	// Read the response.
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
	fmt.Println(body)
	return nil
}


func isFileFound(filePath, token string) bool{
	dat, err := ioutil.ReadFile(filePath)
	check(err)
	sha256 := getSha256(dat)
	
	url := fileURL + sha256 + "/?fields=status"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get() failed with %v", err)
		return false
	}

	if resp.StatusCode == http.StatusNotFound {
		return false
	}

	defer resp.Body.Close()
	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll() failed with: %v", err)
		return false
	}

	var file map[string]interface{}
	if err := json.Unmarshal(jsonBody, &file); err != nil {
		log.Printf("json.Unmarshal() failed with: %v", err)
		return false
	}

	if val, ok := file["status"]; ok {
		status := val.(float64)
		if status == 2 {
			log.Printf("File %s already in DB", sha256)
			return true
		}
	}

	log.Printf("rescanning %s", sha256)
	rescan(sha256, token)
	time.Sleep(15 * time.Second)
	return false
}

func main() {

	if len(os.Args) != 2 {
		log.Println("Usage: uploader <filepath>")
		return
	}
	filePath := os.Args[1]
	_, err := os.Stat(filePath)
    if os.IsNotExist(err) {
		log.Fatalf("%s does not exist", filePath)
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
	filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	// Upload files
	for _, file := range fileList {
		// Check if we have the file already in our database.
		found := isFileFound(file, token)
		if !found {
			upload(file, token)
			time.Sleep(15 * time.Second)
		}
	}
}
