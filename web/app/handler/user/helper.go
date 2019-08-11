package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// generateEmailConfirmationToken creates a JWT token for email confirmation.
func (u *User) generateEmailConfirmationToken() (string, error) {
	rawToken := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := rawToken.Claims.(jwt.MapClaims)
	claims["name"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Generate encoded token and send it as response.
	key := viper.GetString("auth.signkey")
	token, err := rawToken.SignedString([]byte(key))
	return token, err
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
