// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package file

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/labstack/echo"
	"github.com/minio/minio-go"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/common/db"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

type stringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}

// File represent a sample
type File struct {
	Md5       string                 `json:"md5,omitempty"`
	Sha1      string                 `json:"sha1,omitempty"`
	Sha256    string                 `json:"sha256,omitempty"`
	Sha512    string                 `json:"sha512,omitempty"`
	Ssdeep    string                 `json:"ssdeep,omitempty"`
	Crc32     string                 `json:"crc32,omitempty"`
	Magic     string                 `json:"magic,omitempty"`
	Size      int64                  `json:"size,omitempty"`
	Exif      map[string]string      `json:"exif"`
	TriD      []string               `json:"trid"`
	FirstSeen time.Time              `json:"first_seen,omitempty"`
	Strings   []stringStruct         `json:"strings"`
	MultiAV   map[string]interface{} `json:"multiav"`
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
func (file *File) Create() error {
	_, error := db.FilesBucket.Upsert(file.Sha256, file, 0)
	if error != nil {
		log.Fatal(error)
		return error
	}
	log.Infof("File %s added to database.", file.Sha256)
	return nil
}

// GetFileBySHA256 return user document
func GetFileBySHA256(sha256 string) (File, error) {

	// get our file
	file := File{}
	cas, err := db.FilesBucket.Get(sha256, &file)
	if err != nil {
		log.Errorln(err, cas)
		return file, err
	}

	return file, err
}

//=================== /file/sha256 handlers ===================

// GetFile returns file informations.
func GetFile(c echo.Context) error {

	// get path param
	sha256 := c.Param("sha256")

	// ugly
	dir, err := os.Getwd()
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	out := path.Join(dir, "app", "handlers", "file.json")
	raw, err := ioutil.ReadFile(out)
	if err != nil {
		return c.String(http.StatusOK, "something went wrong"+sha256)
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

	// Read the json body
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Validate JSON
	l := gojsonschema.NewBytesLoader(b)
	result, err := app.FileSchema.Validate(l)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !result.Valid() {
		return c.JSON(http.StatusBadRequest, result.Errors())
	}

	// Updates the document.
	file, err := GetFileBySHA256(sha256)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = json.Unmarshal(b, &file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db.FilesBucket.Upsert(sha256, file, 0)
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
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message:     "Missing file",
			Description: "Did you send the file via the form request ?",
		})
	}

	// Check file size
	if fileHeader.Size > app.MaxFileSize {
		return c.JSON(http.StatusRequestEntityTooLarge, Response{
			Message:     "File too large",
			Description: "The maximum allowed is 64MB",
			Filename:    fileHeader.Filename,
		})
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		log.Error("Opening a file handle failed, err: ", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Message:     "Internal error",
			Description: "Internal error",
			Filename:    fileHeader.Filename,
		})
	}
	defer file.Close()

	// Read the content
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("Opening a reading the file content, err: ", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Message:     "Internal error",
			Description: "Internal error",
			Filename:    fileHeader.Filename,
		})
	}

	sha256 := crypto.GetSha256(fileContents)

	// Upload the sample to DO object storage.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	n, err := app.DOClient.PutObjectWithContext(ctx, app.SamplesSpaceBucket,
		sha256, file, fileHeader.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Error("Failed to upload object, err: ", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Message:     "PutObject failed",
			Description: err.Error(),
			Filename:    fileHeader.Filename,
		})
	}
	log.Println("Successfully uploaded bytes: ", n)

	// Save to DB
	NewFile := File{
		Sha256:    sha256,
		FirstSeen: time.Now().UTC(),
		Size:      fileHeader.Size,
	}
	NewFile.Create()

	// Push it to NSQ
	err = app.NsqProducer.Publish("scan", []byte(sha256))
	if err != nil {
		log.Error("Failed to publish to NSQ, err: ", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Message:     "Internal error",
			Description: "Internal error",
			Filename:    fileHeader.Filename,
		})
	}

	// All went fine
	return c.JSON(http.StatusCreated, Response{
		Sha256:      sha256,
		Message:     "ok",
		Description: "File queued successfully for analysis",
		Filename:    fileHeader.Filename,
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
