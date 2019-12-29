// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v6"
	"github.com/saferwall/saferwall/pkg/crypto"
	u "github.com/saferwall/saferwall/pkg/utils"
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/common/db"
	"github.com/saferwall/saferwall/web/app/common/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
	"github.com/xeipuuv/gojsonschema"
)

type stringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}

type submission struct {
	Date     time.Time `json:"date,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Source   string    `json:"source,omitempty"`
	Country  string    `json:"country,omitempty"`
}

// File represent a sample
type File struct {
	Md5             string                 `json:"md5,omitempty"`
	Sha1            string                 `json:"sha1,omitempty"`
	Sha256          string                 `json:"sha256,omitempty"`
	Sha512          string                 `json:"sha512,omitempty"`
	Ssdeep          string                 `json:"ssdeep,omitempty"`
	Crc32           string                 `json:"crc32,omitempty"`
	Magic           string                 `json:"magic,omitempty"`
	Size            int64                  `json:"size,omitempty"`
	Exif            map[string]string      `json:"exif"`
	TriD            []string               `json:"trid"`
	Packer          []string               `json:"packer"`
	FirstSubmission time.Time              `json:"first_submission,omitempty"`
	LastSUbmission  time.Time              `json:"last_submission,omitempty"`
	Submissions     []submission           `json:"submissions,omitempty"`
	Strings         []stringStruct         `json:"strings"`
	MultiAV         map[string]interface{} `json:"multiav"`
	Status          int                    `json:"status"`
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

const (
	queued     = iota
	processing = iota
	finished   = iota
)

// Create creates a new file
func (file *File) Create() error {
	_, error := db.FilesCollection.Upsert(file.Sha256, file, &gocb.UpsertOptions{})
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
	getResult, err := db.FilesCollection.Get(sha256, &gocb.GetOptions{})
	if err != nil {
		log.Errorln(err)
		return file, err
	}

	err = getResult.Content(&file)
	if err != nil {
		log.Errorln(err)
		return file, err
	}
	return file, err
}

// GetAllFiles return all files (optional: selecting fields)
func GetAllFiles(fields []string) ([]File, error) {

	// Select only demanded fields
	var query string
	if len(fields) > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("SELECT ")
		length := len(fields)
		for index, field := range fields {
			buffer.WriteString(field)
			if index < length-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(" FROM `files`")
		query = buffer.String()
	} else {
		query = "SELECT files.* FROM `files`"
	}

	// Execute our query
	rows, err := db.Cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		fmt.Println("Error executing n1ql query:", err)
	}

	// Interfaces for handling streaming return values
	var row File
	var retValues []File

	// Stream the values returned from the query into a typed array of structs
	for rows.Next(&row) {
		retValues = append(retValues, row)
	}

	return retValues, nil
}

// DumpRequest sdsds sd fd
func DumpRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Print(string(requestDump))
	}

}

//=================== /file/sha256 handlers ===================

// GetFile returns file informations.
func GetFile(c echo.Context) error {

	DumpRequest(c.Request())
	// Get user IP
	clientIP := realip.FromRequest(c.Request())
	clientIP2 := utils.GetIPAdress(c.Request())
	clientIP3 := c.RealIP()
	log.Infoln(clientIP)
	log.Infoln(clientIP2)
	log.Infoln(clientIP3)

	// get path param
	sha256 := c.Param("sha256")

	matched, _ := regexp.MatchString("^[a-f0-9]{64}$", sha256)
	if !matched {
		return c.JSON(http.StatusBadRequest, Response{
			Message:     "Invalid sha265",
			Description: "Invalid hash submitted: " + sha256,
		})
	}

	file, err := GetFileBySHA256(sha256)
	if err != nil && gocb.IsKeyNotFoundError(err) {
		return c.JSON(http.StatusNotFound, Response{
			Message:     err.Error(),
			Description: "File was not found in our database",
			Sha256:      sha256,
		})
	}
	return c.JSON(http.StatusOK, file)
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
		for _, desc := range result.Errors() {
			log.Printf("- %s\n", desc)
		}
		return c.JSON(http.StatusBadRequest, errors.New("json validation failed"))
	}

	matched, _ := regexp.MatchString("^[a-f0-9]{64}$", sha256)
	if !matched {
		return c.JSON(http.StatusBadRequest, Response{
			Message:     "Invalid sha265",
			Description: "File hash is not a sha256 hash" + sha256,
			Sha256:      sha256,
		})
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

	db.FilesCollection.Upsert(sha256, file, &gocb.UpsertOptions{})
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
	mgr, err := db.Cluster.Buckets()
	if err != nil {
		log.Errorf("Failed to create bucket manager %v", err)
	}
	err = mgr.FlushBucket("files", nil)
	if err != nil {
		log.Errorf("Failed to flush bucket manager %v", err)
	}
}

// GetFiles returns list of files.
func GetFiles(c echo.Context) error {
	// get query param `fields` for filtering & sanitize them
	filters := utils.GetQueryParamsFields(c)
	if len(filters) > 0 {
		file := File{}
		allowed := utils.IsFilterAllowed(utils.GetStructFields(file), filters)
		if !allowed {
			return c.JSON(http.StatusBadRequest, "Filters not allowed")
		}
	}

	// get all users
	allFiles, err := GetAllFiles(filters)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, allFiles)
}

// PostFiles creates a new file
func PostFiles(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	log.Infoln("New file uploaded by", name)

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

	// Get the size
	size := fileHeader.Size
	log.Infoln("File size: ", size)

	// Read the content
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("Opening a reading the file content, err: ", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Message:     "ReadAll failed",
			Description: "Internal error",
			Filename:    fileHeader.Filename,
		})
	}

	sha256 := crypto.GetSha256(fileContents)
	log.Infoln("File hash: ", sha256)

	// Have we seen this file before
	fileDocument, err := GetFileBySHA256(sha256)
	if err != nil && !gocb.IsKeyNotFoundError(err) {
		return c.JSON(http.StatusInternalServerError, Response{
			Message:     "Something unexpected happened",
			Description: err.Error(),
			Filename:    fileHeader.Filename,
		})
	}

	if gocb.IsKeyNotFoundError(err) {
		// Upload the sample to DO object storage.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		n, err := app.MinioClient.PutObjectWithContext(ctx, app.SamplesSpaceBucket,
			sha256, bytes.NewReader(fileContents), size,
			minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			log.Error("Failed to upload object, err: ", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Message:     "PutObject failed",
				Description: err.Error(),
				Filename:    fileHeader.Filename,
				Sha256:      sha256,
			})
		}
		log.Println("Successfully uploaded bytes: ", n)

		// Save to DB
		now := time.Now().UTC()
		newFile := File{
			Sha256:          sha256,
			FirstSubmission: now,
			LastSUbmission:  now,
			Size:            fileHeader.Size,
			Status:          queued,
		}

		// Get user IP
		clientIP := realip.FromRequest(c.Request())
		clientIP2 := utils.GetIPAdress(c.Request())
		log.Infoln(clientIP)
		log.Infoln(clientIP2)
		ip := net.ParseIP(clientIP)
		country, err := app.GeoIPDB.Country(ip)
		if err != nil {
			log.Error(err)
		}

		// Create new submission
		s := submission{
			Date:     now,
			Filename: fileHeader.Filename,
			Source:   "api",
			Country:  country.Country.IsoCode,
		}
		newFile.Submissions = append(newFile.Submissions, s)
		newFile.Create()

		// Push it to NSQ
		err = app.NsqProducer.Publish("scan", []byte(sha256))
		if err != nil {
			log.Error("Failed to publish to NSQ, err: ", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Message:     "Publish failed",
				Description: "Internal error",
				Filename:    fileHeader.Filename,
				Sha256:      sha256,
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

	// We have already seen this file
	return c.JSON(http.StatusOK, fileDocument)

}

// PutFiles bulk updates of files
func PutFiles(c echo.Context) error {
	return c.String(http.StatusOK, "putFiles")
}

// DeleteFiles delete all files
func DeleteFiles(c echo.Context) error {

	go deleteAllFiles()
	return c.JSON(http.StatusOK, map[string]string{
		"verbose_msg": "ok"})
}

// Download downloads a file.
func Download(c echo.Context) error {
	// get path param
	sha256 := c.Param("sha256")

	reader, err := app.MinioClient.GetObject(
		app.SamplesSpaceBucket, sha256, minio.GetObjectOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer reader.Close()

	_, err = reader.Stat()
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	filepath, err := u.ZipEncrypt(sha256, "infected", reader)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.File(filepath)
}

// Actions over a file. Rescan or Download.
func Actions(c echo.Context) error {

	// Read the json body
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Verify length
	if len(b) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "You have sent an empty json"})
	}

	// Validate JSON
	l := gojsonschema.NewBytesLoader(b)
	result, err := app.FileActionSchema.Validate(l)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if !result.Valid() {
		msg := ""
		for _, desc := range result.Errors() {
			msg += fmt.Sprintf("%s, ", desc.Description())
		}
		msg = strings.TrimSuffix(msg, ", ")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": msg})
	}

	// get the type of action
	var actions map[string]interface{}
	err = json.Unmarshal(b, &actions)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	actionType := actions["type"].(string)

	// get path param
	sha256 := c.Param("sha256")
	matched, _ := regexp.MatchString("^[a-f0-9]{64}$", sha256)
	if !matched {
		return c.JSON(http.StatusBadRequest, Response{
			Message:     "Invalid sha265",
			Description: "File hash is not a sha256 hash" + sha256,
			Sha256:      sha256,
		})
	}

	log.Print(sha256)
	_, err = GetFileBySHA256(sha256)
	if err != nil && gocb.IsKeyNotFoundError(err) {
		return c.JSON(http.StatusNotFound, Response{
			Message:     err.Error(),
			Description: "File was not found in our database",
			Sha256:      sha256,
		})
	}

	if actionType == "rescan" {

		// Push it to NSQ
		err = app.NsqProducer.Publish("scan", []byte(sha256))
		if err != nil {
			log.Error("Failed to publish to NSQ, err: ", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Message:     "Publish failed",
				Description: "Internal error",
				Sha256:      sha256,
			})
		}
		return c.JSON(http.StatusOK, Response{
			Message:     "File rescan successful",
			Description: "Type of action: " + actionType,
			Sha256:      sha256,
		})
	} else if actionType == "download" {
		reader, err := app.MinioClient.GetObject(
			app.SamplesSpaceBucket, sha256, minio.GetObjectOptions{})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer reader.Close()

		_, err = reader.Stat()
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		filepath, err := u.ZipEncrypt(sha256, "infected", reader)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.File(filepath)
	}

	return c.JSON(http.StatusInternalServerError, Response{
		Message:     "Unknown action",
		Description: "Type of action: " + actionType,
		Sha256:      sha256,
	})
}
