package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Authentication permits or rejects access for HTTP requests
type Authentication interface {

	// Authenticate returns ID of identified user creating the HTTP request
	Authenticate(req *http.Request) (userID string, err error)
}

// NotAllowedError is returns when access is not permitted
var NotAllowedError = errors.New("access not permitted")

// UsersAuthentication checks credentials against a list of users
type UsersAuthentication []User

// Authenticate extracts HTTP basic auth user credentials and returns whether a user
// in the list has a matching username and password
func (a UsersAuthentication) Authenticate(req *http.Request) (string, error) {
	name, pass, ok := req.BasicAuth()
	if !ok {
		return "", fmt.Errorf("missing credentials: %w", NotAllowedError)
	}
	for _, user := range a {
		// found a user!
		if user.Name == name && user.Password == pass {
			return user.ID, nil
		}
	}
	return "", NotAllowedError
}

// LoadAuthenticationFromJSON reads a JSON file, returns an Authentication implementation
func LoadAuthenticationFromJSON(filename string) (Authentication, error) {

	// define a slice of users & fill it from a JSON file
	var users []User
	encoded, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(encoded, &users); err != nil {
		return nil, err
	}

	// cast the slice of users into an Authentication implementation
	return UsersAuthentication(users), nil
}
