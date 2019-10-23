package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "calhounio_demo"
	sslmode  = "disable"
)

var (
	// ErrInvalidPort is returned when the configuration provides an invalid port
	ErrInvalidPort = errors.New("invalid port")
)

// PostgresDatabase represents a postgres connection
type PostgresDatabase struct {
	conf PostgresConf
	db   *sql.DB
}

// PostgresConf holds information about the postgres configuration
type PostgresConf struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func parsePostgresConfig(c map[string]string) (PostgresConf, error) {
	pc := PostgresConf{}

	host := "localhost"
	if h, ok := c["host"]; ok {
		host = h
	}

	port := 5432
	if p, ok := c["port"]; ok {
		cp, err := strconv.Atoi(p)
		if err != nil {
			return pc, errors.New("invalid port")
		}
		port = cp
	}

	pc.host = host
	pc.port = port

	return pc, nil
}

// NewPostgresDatabase creates a new postgres connection
func NewPostgresDatabase(opts map[string]string) (*PostgresDatabase, error) {
	pdb := &PostgresDatabase{}

	parsePostgresConfig(opts)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return pdb, nil
}

// Stop closes the database connection
func Stop() {

}
