package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" //required
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Migrate(migrationPath string, db *sqlx.DB, log *zap.Logger) error {

	// Run migrations
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "db: creating driver failed")
	}

	m, merr := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath), "postgres", driver)

	sugar := log.Sugar()
	if merr != nil {
		sugar.Fatalf("migration failed... %v", merr)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		sugar.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	sugar.Info("Database migrated")
	// actual logic to start your application

	return nil
}
