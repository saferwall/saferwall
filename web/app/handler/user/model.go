package user

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/saferwall/saferwall/web/app/common/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/couchbase/gocb.v1"

	"time"
)

var (

	// ErrUserAlreadyConfirmed is retgurned when a user account has been already confirmed.
	ErrUserAlreadyConfirmed = errors.New("Account already confirmed")
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
	log.Infof("User %s was created successefuly", u.Username)
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

// Confirm confirms user.
func Confirm(username string) error {
	user, err := GetByUsername(username)
	if err != nil {
		return err
	}

	if user.Confirmed {
		return ErrUserAlreadyConfirmed
	}

	user.Confirmed = true
	user.Save()
	return nil
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

// DeleteAllUsers will empty users bucket
func DeleteAllUsers() {
	// Keep in mind that you must have flushing enabled in the buckets configuration.
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	err := db.UsersBucket.Manager(username, password).Flush()
	if err != nil {
		log.Print(err)
	}
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

// GetUserByEmail return a user document from email
func GetUserByEmail(email string) (User, error) {

	statement := "SELECT users.* FROM `users` WHERE email=$1"

	// Setup a new query with a placeholder
	query := gocb.NewN1qlQuery(statement)

	// Setup an array for parameters
	var myParams []interface{}
	myParams = append(myParams, email)

	// Interfaces for handling streaming return values
	var row User

	// Execute Query
	rows, err := db.UsersBucket.ExecuteN1qlQuery(query, myParams)
	if err != nil {
		fmt.Println("Error executing n1ql query:", err)
		return row, err
	}

	defer rows.Close()

	err = rows.One(&row)
	if err != nil {
		return row, err
	}

	return row, nil
}

// CreateAdminUser creates admin user.
func CreateAdminUser() {
	username := viper.GetString("app.admin_user")
	password := viper.GetString("app.admin_pwd")
	email := viper.GetString("app.admin_email")

	u, _ := GetByUsername(username)
	if u.Username != "" {
		return
	}

	newUser := User{
		Username: username,
		Email: email,
	}

	t := time.Now().UTC()
	newUser.MemberSince = &t
	newUser.Admin = true
	newUser.Password = hashAndSalt([]byte(password))
	newUser.Confirmed = true
	newUser.Save()
}