package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/server/app/common/utils"
	"github.com/saferwall/saferwall/web/server/app/models/user"
	"github.com/saferwall/saferwall/web/server/app/schemas"
	"github.com/xeipuuv/gojsonschema"
)

// GetUsers get all users
func GetUsers(c echo.Context) error {

	// get query param `fields` for filtering & sanitize them
	filters := utils.GetQueryParamsFields(c)
	if len(filters) > 0 {
		user := user.User{}
		allowed := utils.IsFilterAllowed(user.GetStructFields(), filters)
		if !allowed {
			return c.JSON(http.StatusBadRequest, "Filters not allowed")
		}
	}

	// get all users
	allUsers, err := user.GetAllUsers(filters)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, allUsers)
}

// PostUsers add a new user
func PostUsers(c echo.Context) error {

	/* Read the json body */
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		panic(err.Error())
	}

	// Validate JSON
	documentLoader := gojsonschema.NewBytesLoader(b)
	result, err := schemas.UserSchemaLoader.Validate(documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	// Bind it to our struct
	newUser := user.User{}
	err = json.Unmarshal(b, &newUser)
	if err != nil {
		fmt.Println("Failed to bind json : ", err)
	}

	// Insert to DB
	newUser.MemberSince = time.Now().UTC()
	user.NewUser(newUser)

	return c.JSON(http.StatusCreated, newUser)
}

// PutUsers handles /PUT
func PutUsers(c echo.Context) error {
	return c.String(http.StatusOK, "PutUsers")
}

// DeleteUsers handlers /DELETE
func DeleteUsers(c echo.Context) error {

	// should be processed in the background
	user.DeleteAllUsers()
	return c.JSON(http.StatusOK, "All users deleted")
}

// GetUser handle /GET request
func GetUser(c echo.Context) error {

	// get query param `fields` for filtering & sanitize them
	filters := utils.GetQueryParamsFields(c)
	if len(filters) > 0 {
		user := user.User{}
		allowed := utils.IsFilterAllowed(user.GetStructFields(), filters)
		if !allowed {
			return c.JSON(http.StatusBadRequest, "Filters not allowed")
		}
	}

	// get path param
	username := c.Param("username")
	user, err := user.GetUserByUsernameFields(filters, username)
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

	err := user.DeleteUser(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, username)
}
