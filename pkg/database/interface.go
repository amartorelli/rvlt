package database

import (
	"errors"

	"github.com/amartorelli/rvlt/pkg/model"
)

// Database is the database interface
type Database interface {
	Store(model.User) error
	Get(string) (u model.User, err error)
	Close() error
}

var (
	// ErrInvalidDatabaseType is returned when the database requested is unsupported
	ErrInvalidDatabaseType = errors.New("invalid database type specified")
	// ErrUserNotFound is returned if the user could not be found in the database
	ErrUserNotFound = errors.New("user not found")
)

// NewDatabase returns a new database depending on the type
func NewDatabase(dbType string) (Database, error) {
	switch dbType {
	case "memory":
		db, _ := NewMemoryDatabase()
		return db, nil
	case "postgres":
		db, err := NewPostgresDatabase()
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, ErrInvalidDatabaseType
	}
}
