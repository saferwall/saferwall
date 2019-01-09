package user

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/saferwall/saferwall/web/server/app/common/database"
	"gopkg.in/couchbase/gocb.v1"
)

// User represent a user
type User struct {
	Email       string    `json:"email,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	Confirmed   bool      `json:"confirmed,omitempty"`
	MemberSince time.Time `json:"member_since,omitempty"`
}

// GetStructFields retrieve json struct fields names
func (u User) GetStructFields() []string {

	val := reflect.ValueOf(u)
	var temp string

	var listFields []string
	for i := 0; i < val.Type().NumField(); i++ {
		temp = val.Type().Field(i).Tag.Get("json")
		temp = strings.Replace(temp, ",omitempty", "", -1)
		listFields = append(listFields, temp)
	}

	return listFields
}

func fieldSet(fields ...string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, s := range fields {
		set[s] = true
	}
	return set
}

// SelectFields execlude sensitive fields
func (u *User) SelectFields(fields ...string) map[string]interface{} {
	fs := fieldSet(fields...)
	rt, rv := reflect.TypeOf(*u), reflect.ValueOf(*u)
	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonKey := field.Tag.Get("json")
		if fs[jsonKey] {
			out[jsonKey] = rv.Field(i).Interface()
		}
	}
	return out
}

// NewUser add a user to a database
func NewUser(newUser User) {
	database.UsersBucket.Upsert(newUser.Username, newUser, 0)
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
	rows, err := database.UsersBucket.ExecuteN1qlQuery(query, nil)
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

// DeleteAllUsers will empty users bucket
func DeleteAllUsers() {
	// Keep in mind that you must have flushing enabled in the buckets configuration.
	database.UsersBucket.Manager("", "").Flush()
}

// GetUserByUsername return user document
func GetUserByUsername(username string) (User, error) {

	// get our user
	user := User{}
	cas, err := database.UsersBucket.Get(username, &user)
	if err != nil {
		fmt.Println(err, cas)
		return user, err
	}

	return user, err
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
	rows, err := database.UsersBucket.ExecuteN1qlQuery(query, myParams)
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

// DeleteUser will delete a user
func DeleteUser(username string) error {

	// delete document
	cas, err := database.UsersBucket.Remove(username, 0)
	fmt.Println(cas, err)
	return err
}
