package database

import (
	"os"
	"testing"

	"github.com/amartorelli/rvlt/pkg/model"
)

var db Database

func TestBadDBType(t *testing.T) {
	// test init db
	_, err := NewDatabase("unsupported")
	if err != ErrInvalidDatabaseType {
		t.Error("an unsupported database initialisation should return an error")
	}
}

func TestMemoryStore(t *testing.T) {
	// test storing a user
	err := db.Store(model.User{Username: "john", DOB: "2014-12-03"})
	if err != nil {
		t.Error("storing a user should not fail")
	}
}

func TestMemoryGet(t *testing.T) {
	// test getting a non existent user
	_, err := db.Get("zack")
	if err != ErrUserNotFound {
		t.Error("getting a non existing user should return an error")
	}

	// test getting an existing user
	_ = db.Store(model.User{Username: "john", DOB: "2014-12-03"})

	u, err := db.Get("john")
	if err != nil {
		t.Error("an existing user should be returned")
	}

	if u.Username != "john" || u.DOB != "2014-12-03" {
		t.Error("the user returned does not match the one in the database")
	}
}

func TestMain(m *testing.M) {
	// test init db
	db, _ = NewDatabase("memory")
	os.Exit(m.Run())
}
