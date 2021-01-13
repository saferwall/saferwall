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
	fileURL = "https://api.saferwall.com/v1/files/"
	authURL = "https://api.saferwall.com/v1/auth/login/"
)

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

func rescan(sha256,  authToken string) error {

	log.Printf("rescanning %s\n", sha256)

	payload, err := json.Marshal(map[string]string{
		"type": "rescan",
	})
	check(err)

	url := fileURL + sha256 + "/actions"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	check(err)

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

func scan(sha256 string, authToken string) error {

	log.Printf("Scanning %s\n", sha256)

	url := fileURL + sha256 + "/scan"
	request, err := http.NewRequest("POST", url, nil)
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

func isFileFoundInDB(sha256, token string) bool {

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
	return false
}

func listAllFilesInDb(authToken string) ([]string, error) {

	var listSha256 []string
	url := fileURL + "?fields=sha256"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return listSha256, err
	}

	// Read the response.
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
		return listSha256, err
	}

	var shaMap []map[string]string
	err = json.Unmarshal(body.Bytes(), &shaMap)
	check(err)

	for _, v := range shaMap {
		listSha256 = append(listSha256, v["sha256"])
	}

	resp.Body.Close()
	return listSha256, nil
}
