package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type emailCustomClaims struct {
	Username string `json:"username"`
	Purpose string `json:"purpose"`
	jwt.StandardClaims
}

// generateEmailConfirmationToken creates a JWT token for email confirmation.
func (u *User) generateEmailConfirmationToken() (string, error) {

	// Set custom claims
	claims := &emailCustomClaims{
		u.Username,
		"confirm-email",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	key := viper.GetString("auth.signkey")
	t, err := token.SignedString([]byte(key))
	return t, err
}


// GenerateResetPasswordToken creates a JWT token for password change.
func (u *User) GenerateResetPasswordToken() (string, error) {

	// Set custom claims
	claims := &emailCustomClaims{
		u.Username,
		"confirm-email",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	key := viper.GetString("auth.signkey")
	t, err := token.SignedString([]byte(key))
	return t, err
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
