package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" //required
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres
	"github.com/pkg/errors"
	"log"
	"os"
)

func Migrate(migrationPath string, db *sqlx.DB) error {

	// Run migrations
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "db: creating driver failed")
	}

	m, merr := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath), "postgres", driver)

	if merr != nil {
		log.Fatalf("migration failed... %v", merr)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Println("Database migrated")
	// actual logic to start your application
	os.Exit(0)

	return nil
}

