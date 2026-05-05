package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

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

	if v := os.Getenv("DB_MAX_OPEN_CONNS"); len(v) != 0 {
		if v, err := strconv.Atoi(v); err == nil {
			fmt.Printf("Set sql.DB.MaxOpenConns to %d.\n", v)
			db.SetMaxOpenConns(v)
		} else {
			panic(err)
		}
	}
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); len(v) != 0 {
		if v, err := strconv.Atoi(v); err == nil {
			fmt.Printf("Set sq.DB.SetMaxIdleConns to %d.\n", v)
			db.SetMaxIdleConns(v)
		} else {
			panic(err)
		}
	}
	if v := os.Getenv("DB_CONN_MAX_LIFETIME"); len(v) != 0 {
		if v, err := strconv.Atoi(v); err == nil {
			fmt.Printf("Set sql.DB.SetConnMaxLifetime to %d minutes.\n", v)
			db.SetConnMaxLifetime(time.Duration(v) * time.Minute)
		} else {
			panic(err)
		}
	}
	if v := os.Getenv("DB_CONN_MAX_IDLE_TIME"); len(v) != 0 {
		if v, err := strconv.Atoi(v); err == nil {
			fmt.Printf("Set sql.DB.SetConnMaxIdleTime to %d minutes.\n", v)
			db.SetConnMaxIdleTime(time.Duration(v) * time.Minute)
		} else {
			panic(err)
		}
	}

	return nil
}
