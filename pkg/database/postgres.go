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
	// ErrInvalidConf is returned when the configuration is missing some required fields
	ErrInvalidConf = errors.New("invalid configuration")
)

// PostgresDatabase represents a postgres connection
type PostgresDatabase struct {
	conf    PostgresConf
	connStr string
	db      *sql.DB
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

	var user string
	if u, ok := c["user"]; !ok {
		user = u
	} else {
		return pc, ErrInvalidConf
	}

	var pwd string
	if pw, ok := c["password"]; ok {
		pwd = pw
	} else {
		return pc, ErrInvalidConf
	}

	var dbname string
	if dbn, ok := c["password"]; ok {
		dbname = dbn
	} else {
		return pc, ErrInvalidConf
	}

	pc.host = host
	pc.port = port
	pc.user = user
	pc.password = pwd
	pc.dbname = dbname

	return pc, nil
}

// NewPostgresDatabase creates a new postgres connection
func NewPostgresDatabase(opts map[string]string) (*PostgresDatabase, error) {
	pdb := &PostgresDatabase{}

	conf, err := parsePostgresConfig(opts)
	if err != nil {
		return nil, err
	}
	pdb.conf = conf

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.host, conf.port, conf.user, conf.password, conf.dbname)
	pdb.connStr = psqlInfo

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	pdb.db = db

	fmt.Println("Successfully connected!")
	return pdb, nil
}

// Stop closes the database connection
func (d *PostgresDatabase) Stop() error {
	err := d.db.Close()
	return err
}
