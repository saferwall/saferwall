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
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/common/db"
	"github.com/saferwall/saferwall/web/app/common/utils"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/crypto/bcrypt"
	gocb "gopkg.in/couchbase/gocb.v1"
)

// User represent a user.
type User struct {
	Email       string     `json:"email,omitempty"`
	Username    string     `json:"username,omitempty"`
	Password    string     `json:"password,omitempty"`
	FirstName   string     `json:"first_name,omitempty"`
	LastName    string     `json:"last_name,omitempty"`
	Bio         string     `json:"bio,omitempty"`
	Confirmed   bool       `json:"confirmed,omitempty"`
	MemberSince *time.Time `json:"member_since,omitempty"`
	Admin       bool       `json:"admin,omitempty"`
}

// Save adds user to a database.
func (u *User) Save() {
	db.UsersBucket.Upsert(u.Username, u, 0)
	log.Infof("User was created successefuly: %s", u.Username)
}

// DeleteAllUsers will empty users bucket
func DeleteAllUsers() {
	// Keep in mind that you must have flushing enabled in the buckets configuration.
	db.UsersBucket.Manager("", "").Flush()
}

// GetByUsername return user document
func GetByUsername(username string) (User, error) {

	// get our user
	user := User{}
	cas, err := db.UsersBucket.Get(username, &user)
	if err != nil {
		fmt.Println(err, cas)
		return user, err
	}

	return user, nil
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

// CheckEmailExist returns true if emails exists
func CheckEmailExist(email string) (bool, error) {
	myQuery := "SELECT COUNT(*) as count FROM `users` WHERE email=$1"
	rows, err := db.UsersBucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(myQuery), []interface{}{email})
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var row interface{}
	err = rows.One(&row)
	if err != nil {
		return false, err
	}

	count := row.(map[string]interface{})["count"]
	emailExist := count.(float64) > 0
	return emailExist, nil
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
		allowed := utils.IsFilterAllowed(utils.GetStructFields(user), filters)
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

// Create creates a new user
func Create(username, password, email string) User {

	t := time.Now().UTC()
	u := User{
		Username:    username,
		Password:    password,
		Email:       email,
		MemberSince: &t,
		Admin:       false,
	}

	return u
}

// hashAndSalt hash with a salt a password.
func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
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

	// Verify length
	if len(b) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "You have sent an empty json"})
	}

	// Validate JSON
	l := gojsonschema.NewBytesLoader(b)
	result, err := app.UserSchema.Validate(l)
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

	// Bind it to our User instance.
	newUser := User{}
	err = json.Unmarshal(b, &newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// check if user already exist in DB.
	u, err := GetByUsername(newUser.Username)
	if err == nil && u.Username != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "Username already exists !"})
	}

	// check if email already exists in DB.
	EmailExist, _ := CheckEmailExist(newUser.Email)
	if EmailExist {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "Email already exists !"})
	}

	// Update some details
	t := time.Now().UTC()
	newUser.MemberSince = &t
	newUser.Admin = false
	newUser.Password = hashAndSalt([]byte(newUser.Password))

	// Creates the new user and save it to DB.
	newUser.Save()
	return c.JSON(http.StatusCreated, map[string]string{
		"verbose_msg": "ok"})
}

// GetUsers returns all users.
func GetUsers(c echo.Context) error {

	// get query param `fields` for filtering & sanitize them
	filters := utils.GetQueryParamsFields(c)
	if len(filters) > 0 {
		user := User{}
		allowed := utils.IsFilterAllowed(utils.GetStructFields(user), filters)
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
