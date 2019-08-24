// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package user

import (
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
	"github.com/saferwall/saferwall/web/app/email"
	"github.com/xeipuuv/gojsonschema"
)

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

	// Send confirmation email
	token, err := newUser.generateEmailConfirmationToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"verbose_msg": "Internal server error !"})
	}

	// Generate the email confirmation url
	r := c.Request()
	baseURL := c.Scheme() + "://" + r.Host
	link := baseURL + "/v1/auth/confirm/" + "?token=" + token
	go email.Send(newUser.Username, link, newUser.Email, "confirm")

	return c.JSON(http.StatusCreated, map[string]string{
		"verbose_msg": "ok"})
}

// PutUsers handles /PUT
func PutUsers(c echo.Context) error {
	return c.String(http.StatusOK, "PutUsers")
}

// DeleteUsers handlers /DELETE
func DeleteUsers(c echo.Context) error {

	// should be processed in the background
	go DeleteAllUsers()
	return c.JSON(http.StatusOK, map[string]string{
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
