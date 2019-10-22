package model

import "testing"

func TestIsUsernameValid(t *testing.T) {
	tt := []struct {
		username string
		valid    bool
		msg      string
	}{
		{"john", true, "john should be a valid user"},
		{"John", true, "a user with an uppercase letter should be valid"},
		{"user01", false, "a user containing numbers should not be valid"},
		{"user one", false, "a user with a space should not be valid"},
		{"user!", false, "a user with a special character should not be valid"},
	}

	for _, tc := range tt {
		if isUsernameValid(tc.username) != tc.valid {
			t.Error(tc.msg)
		}
	}
}

func TestIsDOBValid(t *testing.T) {
	tt := []struct {
		dob   string
		valid bool
		msg   string
	}{
		{"2017-02-12", true, "2016-02-12 should be valid"},
		{"2016-02-29", true, "2016 was a leap year and should be valid"},
		{"2015-2-12", false, "a date with the YYYY-M-DD format should be invalid"},
		{"2015-02-1", false, "a date with the YYYY-MM-D format should be invalid"},
		{"3022-02-12", false, "a date in the future should be invalid"},
		{"2015-13-12", false, "a date with a month greather than 12 should be invalid"},
		{"2015-11-43", false, "a date with a day greather than what the month has should be invalid"},
	}

	for _, tc := range tt {
		if isDOBValid(tc.dob) != tc.valid {
			t.Error(tc.msg)
		}
	}
}

func TestIsValid(t *testing.T) {
	tt := []struct {
		u     User
		valid bool
		err   error
		msg   string
	}{
		{User{Username: "john", DOB: "2012-03-02"}, true, nil, "a valid user should be valid"},
		{User{Username: "user01", DOB: "2012-03-02"}, false, ErrInvalidUsername, "a user with a valid DOB but invalid username should be invalid"},
		{User{Username: "user", DOB: "3043-03-02"}, false, ErrInvalidDOB, "a user with a valid username but invalid DOB should be invalid"},
		{User{Username: "user01", DOB: "3043-03-02"}, false, ErrInvalidUsername, "a user with a both invalid username and DOB should be invalid"},
	}

	for _, tc := range tt {
		v, err := tc.u.IsValid()
		if v != tc.valid || err != tc.err {
			t.Error(tc.msg)
		}
	}
}
