package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

const (
	up   = "up"
	down = "down"
)

func main() {
	direction := up
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}
	if direction != up && direction != down {
		panic(fmt.Sprintf("Invalid argument: %s", direction))
	}

	var err error
	ptn := fmt.Sprintf("cmd/migrate/migrations/*.%s.sql", direction)
	files, err := filepath.Glob(ptn)
	if err != nil {
		panic(err)
	} else if len(files) == 0 {
		fmt.Println("No migration files found matching: ", ptn)
		return
	}
	if direction == up {
		// 001_*.up.sql -> 005_*.up.sql
		sort.Sort(sort.StringSlice(files))
	} else {
		// 005_*.down.sql -> 001_*.down.sql
		sort.Sort(sort.Reverse(sort.StringSlice(files)))
	}

	// DSN is implicitly constructed from env variables such as PGHOST.
	db, err := sql.Open("postgres", "sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		panic(err)
	}

	// e.g., "001", "020", "999"
	vPtn := regexp.MustCompile(`^\d{3}$`)
	for _, f := range files {
		// "001_create_table.up.sql" -> "001"
		v := strings.SplitN(filepath.Base(f), "_", 2)[0]
		if !vPtn.MatchString(v) {
			panic(fmt.Sprint("Invalid migration file name: ", f))
		}

		var rslt string
		err = db.QueryRow(
			`SELECT version FROM schema_migrations WHERE version = $1`, v,
		).Scan(&rslt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Failed to check migration record for: ", f)
			panic(err)
		} else if (direction == up && err == nil) ||
			(direction == down && err != nil) {
			fmt.Println("Skipping (already applied): ", f)
			continue
		}

		migration, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("Can not read ", f)
			panic(err)
		}

		tx, err := db.Begin()
		if err != nil {
			fmt.Println("Failed to begin transaction for: ", f)
			panic(err)
		}

		_, err = tx.Exec(string(migration))
		if err != nil {
			tx.Rollback()
			fmt.Println("Failed to apply: ", f)
			panic(err)
		}

		if direction == up {
			_, err = tx.Exec(`INSERT INTO schema_migrations(version) VALUES($1)`, v)
		} else {
			_, err = tx.Exec(`DELETE FROM schema_migrations WHERE version = $1`, v)
		}
		if err != nil {
			tx.Rollback()
			fmt.Println("Failed to update migration record for: ", f)
			panic(err)
		}

		if err := tx.Commit(); err != nil {
			fmt.Println("Failed to commit: ", f)
			panic(err)
		}

		if direction == up {
			fmt.Println("Applied: ", f)
		} else {
			fmt.Println("Reverted: ", f)
		}
	}
}
