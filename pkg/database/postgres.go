package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/amartorelli/rvlt/pkg/model"
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
	sslmode  string
}

func parsePostgresConfig() (PostgresConf, error) {
	pc := PostgresConf{}

	host := "localhost"
	if h := os.Getenv("POSTGRES_HOST"); h != "" {
		host = h
	}

	port := 5432
	if p := os.Getenv("POSTGRES_PORT"); p != "" {
		cp, err := strconv.Atoi(p)
		if err != nil {
			return pc, errors.New("invalid port")
		}
		port = cp
	}

	var user string
	if u := os.Getenv("POSTGRES_USER"); u != "" {
		user = u
	} else {
		return pc, ErrInvalidConf
	}

	var pwd string
	if pw := os.Getenv("POSTGRES_PASSWORD"); pw != "" {
		pwd = pw
	} else {
		return pc, ErrInvalidConf
	}

	var dbname string
	if dbn := os.Getenv("POSTGRES_DB"); dbn != "" {
		dbname = dbn
	} else {
		return pc, ErrInvalidConf
	}

	var sslmode string = "disable"
	if sslm := os.Getenv("POSTGRES_SSLMODE"); sslm == "enable" || sslm == "disable" {
		sslmode = sslm
	} else {
		return pc, ErrInvalidConf
	}

	pc.host = host
	pc.port = port
	pc.user = user
	pc.password = pwd
	pc.dbname = dbname
	pc.sslmode = sslmode

	return pc, nil
}

// NewPostgresDatabase creates a new postgres connection
func NewPostgresDatabase() (*PostgresDatabase, error) {
	pdb := &PostgresDatabase{}

	conf, err := parsePostgresConfig()
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

var (
	querySelectUser = `SELECT * FROM birthdays WHERE username = '$1'`
	queryInsertUser = `INSERT INTO birthdays(username, birthday) VALUES('$1', '$2')`
	queryUpdateUser = `UPDATE birthdays SET birthday = $1 WHERE username = '$2'`
)

// Store stores a user in postgres
func (d *PostgresDatabase) Store(u model.User) error {
	_, err := d.Get(u.Username)
	// if the user is not present we insert
	if err == ErrUserNotFound {
		_, err := d.db.Query(queryInsertUser, u.Username, u.DOB)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// if we got here, the user is already present and we should do an update
	stmt, err := d.db.Prepare(queryUpdateUser)
	_, err = stmt.Exec(queryUpdateUser, u.DOB, u.Username)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves a user's birthday from postgres
func (d *PostgresDatabase) Get(user string) (u model.User, err error) {
	usr := model.User{}

	rows, err := d.db.Query(querySelectUser, user)
	if err != nil {
		return usr, err
	}

	if !rows.Next() {
		return usr, ErrUserNotFound
	}

	err = rows.Scan(&usr.Username, &usr.DOB)
	if err != nil {
		return usr, err
	}

	return usr, nil
}

// Close closes the connection to the DB
func (d *PostgresDatabase) Close() error {
	return d.db.Close()
}
