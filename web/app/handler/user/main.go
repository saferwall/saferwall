// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/common/db"
	"github.com/saferwall/saferwall/web/app/common/utils"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	gocb "gopkg.in/couchbase/gocb.v1"
)

// User represent a user.
type User struct {
	Email       string    `json:"email,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	Confirmed   bool      `json:"confirmed,omitempty"`
	MemberSince time.Time `json:"member_since,omitempty"`
}

// GetStructFields retrieve json struct fields names
func (u User) GetStructFields() []string {

	val := reflect.ValueOf(u)
	var temp string

	var listFields []string
	for i := 0; i < val.Type().NumField(); i++ {
		temp = val.Type().Field(i).Tag.Get("json")
		temp = strings.Replace(temp, ",omitempty", "", -1)
		listFields = append(listFields, temp)
	}

	return listFields
}

func fieldSet(fields ...string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, s := range fields {
		set[s] = true
	}
	return set
}

// SelectFields execlude sensitive fields
func (u *User) SelectFields(fields ...string) map[string]interface{} {
	fs := fieldSet(fields...)
	rt, rv := reflect.TypeOf(*u), reflect.ValueOf(*u)
	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonKey := field.Tag.Get("json")
		if fs[jsonKey] {
			out[jsonKey] = rv.Field(i).Interface()
		}
	}
	return out
}

// Create add a new user to a database.
func (u *User) Create() {

	u.MemberSince = time.Now().UTC()
	db.UsersBucket.Upsert(u.Username, u, 0)
	log.Infof("User was created successefuly: %s", u.Username)
}

// DeleteAllUsers will empty users bucket
func DeleteAllUsers() {
	// Keep in mind that you must have flushing enabled in the buckets configuration.
	db.UsersBucket.Manager("", "").Flush()
}

// GetUserByUsername return user document
func GetUserByUsername(username string) (User, error) {

	// get our user
	user := User{}
	cas, err := db.UsersBucket.Get(username, &user)
	if err != nil {
		fmt.Println(err, cas)
		return user, err
	}

	return user, err
}

// GetUserByUsernameFields return user by username(optional: selecting fields)
func GetUserByUsernameFields(fields []string, username string) (User, error) {

	// Select only demanded fields
	var statement string
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
		buffer.WriteString(" FROM `users` WHERE username=$1")
		statement = buffer.String()
	} else {
		statement = "SELECT users.* FROM `users` WHERE username=$1"
	}

	// Setup a new query with a placeholder
	query := gocb.NewN1qlQuery(statement)

	// Setup an array for parameters
	var myParams []interface{}
	myParams = append(myParams, username)

	// Interfaces for handling streaming return values
	var row User

	// Execute Query
	rows, err := db.UsersBucket.ExecuteN1qlQuery(query, myParams)
	if err != nil {
		fmt.Println("Error executing n1ql query:", err)
		return row, err
	}

	// Stream the first result only into the interface
	err = rows.One(&row)
	if err != nil {
		fmt.Println("ERROR ITERATING QUERY RESULTS:", err)
		return row, err
	}

	return row, nil
}

// deleteUser will delete a user
func deleteUser(username string) error {

	// delete document
	cas, err := db.UsersBucket.Remove(username, 0)
	fmt.Println(cas, err)
	return err
}

// GetUser handle /GET request
func GetUser(c echo.Context) error {

	// get query param `fields` for filtering & sanitize them
	filters := utils.GetQueryParamsFields(c)
	if len(filters) > 0 {
		user := User{}
		allowed := utils.IsFilterAllowed(user.GetStructFields(), filters)
		if !allowed {
			return c.JSON(http.StatusBadRequest, "Filters not allowed")
		}
	}

	// get path param
	username := c.Param("username")
	user, err := GetUserByUsernameFields(filters, username)
	if err != nil {
		return c.JSON(http.StatusNotFound, username)
	}
	return c.JSON(http.StatusOK, user)
}

// PostUser handle /POST request
func PostUser(c echo.Context) error {
	return c.String(http.StatusOK, "PostUser")
}

// PutUser handle /PUT request
func PutUser(c echo.Context) error {
	return c.String(http.StatusOK, "PutUser")
}

// DeleteUser handle /DELETE request
func DeleteUser(c echo.Context) error {

	// get path param
	username := c.Param("username")

	err := deleteUser(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, username)
}

// GetAllUsers return all users (optional: selecting fields)
func GetAllUsers(fields []string) ([]User, error) {

	// Select only demanded fields
	var statement string
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
		buffer.WriteString(" FROM `users`")
		statement = buffer.String()
	} else {
		statement = "SELECT users.* FROM `users`"
	}

	// Execute our query
	query := gocb.NewN1qlQuery(statement)
	rows, err := db.UsersBucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		fmt.Println("Error executing n1ql query:", err)
	}

	// Interfaces for handling streaming return values
	var row User
	var retValues []User

	// Stream the values returned from the query into a typed array of structs
	for rows.Next(&row) {
		retValues = append(retValues, row)
	}

	return retValues, nil
}

// PostUsers adds a new user.
func PostUsers(c echo.Context) error {

	// Read the json body
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate JSON
	l := gojsonschema.NewBytesLoader(b)
	result, err := app.UserSchemaLoader.Validate(l)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !result.Valid() {
		return c.JSON(http.StatusBadRequest, result.Errors())
	}

	// Bind it to our User instance.
	usr := User{}
	err = json.Unmarshal(b, &usr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Creates the new user and save it to DB.
	usr.Create()
	return c.JSON(http.StatusCreated, usr)
}

// GetUsers returns all users.
func GetUsers(c echo.Context) error {

	// get query param `fields` for filtering & sanitize them
	filters := utils.GetQueryParamsFields(c)
	if len(filters) > 0 {
		user := User{}
		allowed := utils.IsFilterAllowed(user.GetStructFields(), filters)
		if !allowed {
			return c.JSON(http.StatusBadRequest, "Filters not allowed")
		}
	}

	// get all users
	allUsers, err := GetAllUsers(filters)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, allUsers)
}

// PutUsers handles /PUT
func PutUsers(c echo.Context) error {
	return c.String(http.StatusOK, "PutUsers")
}

// DeleteUsers handlers /DELETE
func DeleteUsers(c echo.Context) error {

	// should be processed in the background
	DeleteAllUsers()
	return c.JSON(http.StatusOK, "All users deleted")
}
