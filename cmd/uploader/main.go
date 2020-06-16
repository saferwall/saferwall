// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	fileURL = "https://api.saferwall.com/v1/files/"
	authURL = "https://api.saferwall.com/v1/auth/login/"
	bucket  = "saferwall-samples"
	region  = "us-east-1"
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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func difference(a, b []string) []string {
    mb := make(map[string]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []string
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
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

	log.Printf("rescanning %s\n", sha256)

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

func isFileFoundInObjStorage(sha256 string) bool {

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create S3 service client
	svc := s3.New(sess)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(sha256)})
	if err != nil {
		log.Printf("svc.HeadObject failed with: %v", err)
		return false
	}

	return true
}

func listObject(bucket, region string, verbose bool) []string {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	check(err)

	// Create S3 service client
	svc := s3.New(sess)

	// Get the list of items
	var objKeys []string
	query := &s3.ListObjectsV2Input{
		Bucket: &bucket,
	}

	for {
		resp, err := svc.ListObjectsV2(query)
		if err != nil {
			exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
		}

		for _, item := range resp.Contents {
			if verbose {
				fmt.Println("Name:         ", *item.Key)
				fmt.Println("Last modified:", *item.LastModified)
				fmt.Println("Size:         ", *item.Size)
				fmt.Println("Storage class:", *item.StorageClass)
				fmt.Println("")
			}

			objKeys = append(objKeys, *item.Key)
		}

		// Fetch the next chunk.
		query.ContinuationToken = resp.NextContinuationToken

		if *resp.IsTruncated == false {
			break
		}
	}

	fmt.Println("Found", len(objKeys), "items in bucket", bucket)
	fmt.Println("")
	return objKeys
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

	shaListInDB, err := listAllFilesInDb(token)
	check(err)
	shaListInObjStorage := listObject(bucket, region, false)

	missingSha := difference(shaListInObjStorage, shaListInDB)
	fmt.Print(missingSha)

	// Walk over directory.
	fileList := []string{}
	filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	// Upload files
	for _, filename := range fileList {

		// Get sha256
		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Failed to read %s, reason: %v", filePath, err)
		}
		sha256 := getSha256(dat)

		// Check if we have the file already in our database.
		found := isFileFoundInDB(sha256, token)
		if !found {
			if isFileFoundInObjStorage(sha256) {
				scan(sha256, token)
			} else {
				upload(filename, token)
			}
		} else {
			rescan(sha256, token)
		}
		time.Sleep(15 * time.Second)
	}
}
