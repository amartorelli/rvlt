package database

import "github.com/amartorelli/rvlt/pkg/model"

// MemoryDatabase is an in-memory map to mock the database
type MemoryDatabase struct {
	users map[string]string
}

// NewMemoryDatabase creates a new MemoryDatabase
func NewMemoryDatabase() (*MemoryDatabase, error) {
	return &MemoryDatabase{}, nil
}

// Store stores a user in memory
func (d *MemoryDatabase) Store(u model.User) error {
	d.users[u.Username] = u.DOB
	return nil
}

// Get retrieves a user's birthday from memory
func (d *MemoryDatabase) Get(user string) (u model.User, err error) {
	usr := model.User{}
	dob, ok := d.users[user]
	if !ok {
		return usr, ErrUserNotFound
	}

	usr.Username = user
	usr.DOB = dob

	return usr, nil
}
