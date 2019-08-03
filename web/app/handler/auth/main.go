package auth

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/saferwall/saferwall/web/app/handler/user"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

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

// Register handle new user sign-up
func Register(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Create a new instance & validate input
	newUser, err := user.Create(username, password, email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": err.Error()})
	}

	// check if user already exist in DB.
	u, err := user.GetByUsername(username)
	if err == nil && u.Username != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "Username already exists !"})
	}

	// check if email already exists in DB.
	EmailExist, _ := user.CheckEmailExist(email)
	if EmailExist {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": "Email already exists !"})
	}

	// Store the password hash instead.
	newUser.Password = hashAndSalt([]byte(password))

	newUser.Save()
	return c.JSON(http.StatusCreated, map[string]string{
		"verbose_msg": "ok",
	})
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

func validateLoginInput(username, password string) error {
	r := regexp.MustCompile(`^[a-zA-Z0-9]{1,20}$`)
	if !r.MatchString(username) {
		return errors.New("Username should be alpha-numeric between 1 and 20 length char")
	}

	if password == "" {
		return errors.New("The password field is required")
	}
	return nil
}

// Login handle user authentication
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Validate input format
	err := validateLoginInput(username, password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"verbose_msg": err.Error()})
	}

	usr, err := user.GetByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"verbose_msg": "Username does not exist !"})
	}

	if !comparePasswords(usr.Password, []byte(password)) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"verbose_msg": "Username or password does not match !"})
	}

	token, err := createJwtToken(usr)
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
