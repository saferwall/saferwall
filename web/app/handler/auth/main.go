package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/saferwall/saferwall/web/app"
	"github.com/saferwall/saferwall/web/app/middleware"
	"github.com/saferwall/saferwall/web/app/email"
	"github.com/saferwall/saferwall/web/app/handler/user"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/crypto/bcrypt"
)

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// IsAdmin returns true if user is admin
func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["admin"].(bool)
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

// createJwtToken creates a JWT token.
func createJwtToken(u user.User) (string, error) {
	// Set custom claims
	claims := &middleware.LoginCustomClaims{
		u.Username,
		false,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	key := viper.GetString("auth.signkey")
	t, err := token.SignedString([]byte(key))
	return t, err

}

// create cookie to hold the JWT token.
func createJwtCookie(token string) *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = "JWTCookie"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 72)
	cookie.Path = "/"
	cookie.HttpOnly = false 
	cookie.Secure = false  

	// if app.Debug {
	// 	cookie.HttpOnly = false 
	// 	cookie.Secure = false   
	// } else {
	// 	cookie.HttpOnly = true 
	// 	cookie.Secure = true   
	// }

	// cookie.SameSite = http.SameSiteLaxMode
	// cookie.Domain = "api.saferwall.com"
	return cookie
}

// Login handle user authentication
func Login(c echo.Context) error {
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
	result, err := app.LoginSchema.Validate(l)
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
	loginUser := user.User{}
	err = json.Unmarshal(b, &loginUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	loginUsername := loginUser.Username
	loginPassword := loginUser.Password
	u, err := user.GetByUsername(loginUsername)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"verbose_msg": "Username does not exist !"})
	}

	if !comparePasswords(u.Password, []byte(loginPassword)) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"verbose_msg": "Username or password does not match !"})
	}

	if !u.Confirmed {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"verbose_msg": "Account not confirmed, please confirm your email !"})
	}

	token, err := createJwtToken(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"verbose_msg": "Internal server error !"})
	}

	// Create a cookie to place the jwt token
	cookie := createJwtCookie(token)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"verbose_msg": "You were logged in !",
		"token":       token,
	})
}


// ResetPassword handle password reset.
func ResetPassword(c echo.Context) error {
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
	result, err := app.EmailSchema.Validate(l)
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
	var jsonEmail map[string]interface{}
	err = json.Unmarshal(b, &jsonEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	resetEmail := jsonEmail["email"].(string)

	// check if email already exists in DB.
	EmailExist, _ := user.CheckEmailExist(resetEmail)
	if !EmailExist {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "Email does not exists !"})
	}

	u, err := user.GetUserByEmail(resetEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	token, err := u.GenerateResetPasswordToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Generate the email confirmation url
	link := viper.GetString("ui.address") + "/reset-password" + "?token=" + token
	go email.Send(u.Username, link, u.Email, "reset")

	return c.JSON(http.StatusOK, map[string]string{
		"verbose_msg": "ok",
	})
}

// Admin shows admin
func Admin(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{
		"verbose_msg": "You are admin"})
}

// Confirm confirms a user account.
func Confirm(c echo.Context) error {

	// get path param
	token := c.QueryParam("token")

	if token == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"verbose_msg": "You provided an empty token!"})
	}

	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*middleware.CustomClaims)
	tokenType := claims.Purpose
	if tokenType != "confirm-email" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "You provided an invalid token type!"})
	}

	username := claims.Username

	// Confirm account
	err := user.Confirm(username)
	if err != nil {
		if err == user.ErrUserAlreadyConfirmed {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"verbose_msg": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"verbose_msg": "Internal server error !"})
	}

	url := viper.GetString("ui.address") + "/login"
	return c.Redirect(http.StatusPermanentRedirect, url)
}


// NewPassword handle password reset.
func NewPassword(c echo.Context) error {

	// get path param
	token := c.QueryParam("token")

	if token == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"verbose_msg": "You provided an empty token!"})
	}

	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*middleware.CustomClaims)
	tokenType := claims.Purpose
	if tokenType != "reset-password" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "You provided an invalid token type!"})
	}

	usr, err := user.GetByUsername(claims.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"verbose_msg": "Username does not exist !"})
	}

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
	result, err := app.ResetPasswordSchema.Validate(l)
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
	var jsonPassword map[string]interface{}
	err = json.Unmarshal(b, &jsonPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	newPassword := jsonPassword["password"].(string)
	usr.UpdatePassword(newPassword)

	return c.JSON(http.StatusOK, map[string]string{
		"verbose_msg": "ok",
	})
}
