package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/amartorelli/rvlt/pkg/model"
	log "github.com/sirupsen/logrus"

	// used to psql
	_ "github.com/lib/pq"
)

var (
	// ErrInvalidPort is returned when the configuration provides an invalid port
	ErrInvalidPort = errors.New("invalid port")
	// ErrInvalidUser is returned when the configuration provides an invalid user
	ErrInvalidUser = errors.New("invalid user")
	// ErrInvalidPassword is returned when the configuration provides an invalid password
	ErrInvalidPassword = errors.New("invalid password")
	// ErrInvalidDatabase is returned when the configuration provides an invalid database
	ErrInvalidDatabase = errors.New("invalid database")
	// ErrInvalidSSLMode is returned when the configuration provides an invalid ssl mode
	ErrInvalidSSLMode = errors.New("invalid ssl mode")
)

var (
	querySelectUser = "SELECT * FROM birthdays WHERE username = $1"
	queryInsertUser = "INSERT INTO birthdays(username, birthday) VALUES($1, $2)"
	queryUpdateUser = "UPDATE birthdays SET birthday = $1 WHERE username = $2"
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
			return pc, ErrInvalidPort
		}
		port = cp
	}

	var user string
	if u := os.Getenv("POSTGRES_USER"); u != "" {
		user = u
	} else {
		return pc, ErrInvalidUser
	}

	var pwd string
	if pw := os.Getenv("POSTGRES_PASSWORD"); pw != "" {
		pwd = pw
	} else {
		return pc, ErrInvalidPassword
	}

	var dbname string
	if dbn := os.Getenv("POSTGRES_DB"); dbn != "" {
		dbname = dbn
	} else {
		return pc, ErrInvalidDatabase
	}

	var sslmode string = "disable"
	if sslm := os.Getenv("POSTGRES_SSLMODE"); sslm == "enable" || sslm == "disable" {
		sslmode = sslm
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

	log.Info("db connection initialised")
	return pdb, nil
}

// Stop closes the database connection
func (d *PostgresDatabase) Stop() error {
	err := d.db.Close()
	return err
}

// Store stores a user in postgres
func (d *PostgresDatabase) Store(u model.User) error {
	opMetric.WithLabelValues("store-user").Inc()

	ostart := time.Now()
	defer opDuration.WithLabelValues("store-user").Observe(time.Since(ostart).Seconds())

	dob, err := time.Parse("2006-01-02", u.DOB)
	if err != nil {
		opErrMetric.WithLabelValues("store-user").Inc()
		return err
	}
	_, err = d.Get(u.Username)
	// if the user is not present we insert
	if err == ErrUserNotFound {
		_, err := d.db.Query(queryInsertUser, u.Username, dob)
		if err != nil {
			opErrMetric.WithLabelValues("store-user").Inc()
			return err
		}
		opMetric.WithLabelValues("store-user").Inc()
		return nil
	} else if err != nil {
		opErrMetric.WithLabelValues("store-user").Inc()
		return err
	}

	// if we got here, the user is already present and we should do an update
	stmt, err := d.db.Prepare(queryUpdateUser)
	_, err = stmt.Exec(dob, u.Username)
	if err != nil {
		opErrMetric.WithLabelValues("store-user").Inc()
		return err
	}

	return nil
}

// Get retrieves a user's birthday from postgres
func (d *PostgresDatabase) Get(user string) (u model.User, err error) {
	opMetric.WithLabelValues("get-user").Inc()

	ostart := time.Now()
	defer opDuration.WithLabelValues("get-user").Observe(time.Since(ostart).Seconds())

	usr := model.User{}

	rows, err := d.db.Query(querySelectUser, user)
	if err != nil {
		opErrMetric.WithLabelValues("get-user").Inc()
		return usr, err
	}

	if !rows.Next() {
		opErrMetric.WithLabelValues("get-user").Inc()
		return usr, ErrUserNotFound
	}

	var dob time.Time
	err = rows.Scan(&usr.Username, &dob)
	if err != nil {
		opErrMetric.WithLabelValues("get-user").Inc()
		return usr, err
	}
	usr.DOB = dob.Format("2006-01-02")

	return usr, nil
}

// Close closes the connection to the DB
func (d *PostgresDatabase) Close() error {
	return d.db.Close()
}
