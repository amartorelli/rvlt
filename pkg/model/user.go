package model

import (
	"errors"
	"regexp"
	"time"
)

// User represents a user
type User struct {
	Username string
	DOB      string
}

var (
	// ErrInvalidUsername is retuned when the username provided is invalid
	ErrInvalidUsername = errors.New("invalid username")
	// ErrInvalidDOB is retuned when the date of birth provided is invalid
	ErrInvalidDOB = errors.New("invalid date of birth")
)

var usernameRe = regexp.MustCompile("^[a-zA-Z]+$")
var dobRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// IsValid returns true if the user's info respect the rules
func (u User) IsValid() (bool, error) {
	if !isUsernameValid(u.Username) {
		return false, ErrInvalidUsername
	}

	if !isDOBValid(u.DOB) {
		return false, ErrInvalidDOB
	}

	return true, nil
}

func isUsernameValid(u string) bool {
	return usernameRe.MatchString(u)
}

func isDOBValid(dob string) bool {
	// return false if it doesn't parse as a YYYY-MM-DD date
	t, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return false
	}

	// return false if DOB is in the future
	if t.After(time.Now()) {
		return false
	}

	return true
}
