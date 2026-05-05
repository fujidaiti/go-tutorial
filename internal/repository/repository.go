package repository

import (
	"database/sql"

	// Init the PostgresSQL driver.
	_ "github.com/lib/pq"
)

var db *sql.DB

// Db returns a singleton database client.
// Make sure to call Init first.
func Db() *sql.DB {
	return db
}

// Init initializes the repository instance ensureing that the database has opend
// a connection was successfuly established.
//
// Make sure to call `defer repository.Db.Close()` if Init succeeds.
func Init() error {
	if db != nil {
		panic("Do not call Init multiple times.")
	}

	var err error
	// DSN is implicitly constructed from env variables such as PGHOST.
	db, err = sql.Open("postgres", "sslmode=disable")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
