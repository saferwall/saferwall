// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package file

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/common/db"
)

// File represent a sample
type File struct {
	Md5       string    `json:"md5,omitempty"`
	Sha1      string    `json:"sha1,omitempty"`
	Sha256    string    `json:"sha256,omitempty"`
	Sha512    string    `json:"sha512,omitempty"`
	Ssdeep    string    `json:"ssdeep,omitempty"`
	Crc32     string    `json:"crc32,omitempty"`
	Magic     string    `json:"magic,omitempty"`
	Size      int       `json:"size,omitempty"`
	FirstSeen time.Time `json:"first_seen,omitempty"`
}

// Response JSON
type Response struct {
	Sha256      string `json:"sha256,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
	Filename    string `json:"filename,omitempty"`
}

// AV vendor
type AV struct {
	Vendor string `json:"vendor,omitempty"`
}

// Create creates a new file
func (file *File) Create(documentID string) error {
	_, error := db.FilesBucket.Upsert(documentID, file, 0)
	if error != nil {
		log.Fatal(error)
		return error
	}

	return nil
}

// GetFileBySHA256 return user document
func GetFileBySHA256(sha256 string) (File, error) {

	// get our file
	file := File{}
	cas, err := db.FilesBucket.Get(sha256, &file)
	if err != nil {
		fmt.Println(err, cas)
		return file, err
	}

	return file, err
}

//=================== /file/sha256 handlers ===================

// GetFile returns file informations.
func GetFile(c echo.Context) error {

	// get path param
	// sha256 := c.Param("sha256")

	// ugly
	dir, err := os.Getwd()
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	out := path.Join(dir, "app", "handlers", "file.json")
	raw, err := ioutil.ReadFile(out)
	if err != nil {
		return c.String(http.StatusOK, "something went wrong")
	}
	// r := Response{Sha256: sha256}
	var my map[string]interface{}
	json.Unmarshal(raw, &my)
	return c.JSON(http.StatusOK, my)
}

// PutFile updates a specific file
func PutFile(c echo.Context) error {

	// get path param
	sha256 := c.Param("sha256")

	return c.JSON(http.StatusOK, sha256)
}

// DeleteFile deletes a specific file
func DeleteFile(c echo.Context) error {

	// get path param
	sha256 := c.Param("sha256")
	return c.JSON(http.StatusOK, sha256)
}

// deleteAllFiles will empty files bucket
func deleteAllFiles() {
	// Keep in mind that you must have flushing enabled in the buckets configuration.
	db.FilesBucket.Manager("", "").Flush()
}

// GetFiles returns list of files.
func GetFiles(c echo.Context) error {
	return c.String(http.StatusOK, "getFiles")
}

// PostFiles creates a new file
func PostFiles(c echo.Context) error {

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message:     "Missing file",
			Description: "Did you send the file via the form request ?",
			Filename:    file.Filename,
		})
	}

	// Check file size
	if file.Size > app.MaxFileSize {
		return c.JSON(http.StatusRequestEntityTooLarge, Response{
			Message:     "File too large",
			Description: "The maximum allowed is 64MB",
			Filename:    file.Filename,
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read the content
	fileContents, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	// Destination
	FilePath := path.Join(app.StoragePath, file.Filename)
	dst, err := os.Create(FilePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Writefile
	ioutil.WriteFile(FilePath, fileContents, 0644)
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Save to DB
	Sha256 := crypto.GetSha256(fileContents)
	NewFile := File{
		Sha256:    Sha256,
		FirstSeen: time.Now().UTC(),
	}
	NewFile.Create(Sha256)

	// Push it to NSQ

	return c.JSON(http.StatusCreated, Response{
		Sha256:      Sha256,
		Message:     "ok",
		Description: "File queued successfully for analysis",
		Filename:    file.Filename,
	})
}

// PutFiles bulk updates of files
func PutFiles(c echo.Context) error {
	return c.String(http.StatusOK, "putFiles")
}

// DeleteFiles delete all files
func DeleteFiles(c echo.Context) error {

	deleteAllFiles()
	return c.String(http.StatusOK, "deleteFiles")
}
