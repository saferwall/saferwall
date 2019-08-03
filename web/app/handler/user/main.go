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
	"regexp"
	"time"

	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/common/db"
	"github.com/saferwall/saferwall/web/app/common/utils"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	gocb "gopkg.in/couchbase/gocb.v1"
	"gopkg.in/go-playground/validator.v9"
)

// User represent a user.
type User struct {
	Email       string     `json:"email,omitempty" validate:"required,email"`
	Username    string     `json:"username,omitempty" validate:"required,username"`
	Password    string     `json:"password,omitempty" validate:"required,min=8"`
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
func Create(username, password, email string) (User, error) {

	t := time.Now().UTC()
	u := User{
		Username:    username,
		Password:    password,
		Email:       email,
		MemberSince: &t,
		Admin:       false,
	}

	err := validate(u)
	return u, err
}

// validate user input during account creation
func validate(u User) error {
	v := validator.New()
	var r = regexp.MustCompile(`^[a-zA-Z0-9]{1,20}$`)
	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return r.MatchString(fl.Field().String())
	})

	err := v.Struct(u)
	return err
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
	result, err := app.UserSchema.Validate(l)
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
	usr.Save()
	return c.JSON(http.StatusCreated, usr)
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
