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
	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/app"
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

	rawToken := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := rawToken.Claims.(jwt.MapClaims)
	claims["name"] = u.Username
	claims["admin"] = u.Admin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	key := viper.GetString("auth.signkey")
	token, err := rawToken.SignedString([]byte(key))
	return token, err
}

// create cookie to hold the JWT token.
func createJwtCookie(token string) *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = "JWTCookie"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 72)
	cookie.Path = "/"
	cookie.HttpOnly = false // change this later
	cookie.Secure = false   // change this later
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
	claims := u.Claims.(jwt.MapClaims)
	username := claims["name"].(string)

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

	return c.Redirect(http.StatusPermanentRedirect, "http://localhost:8081/#/upload")
	// return c.JSON(http.StatusAccepted, map[string]string{
	// 	"verbose_msg": "Account confirmed !"})
}

