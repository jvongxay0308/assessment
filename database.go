package main

import (
	// imported to register the postgres migration driver
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// imported to register the file source migration driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// imported to register the postgres database driver
	_ "github.com/lib/pq"
)

// tryToMigrate attempts to apply any pending migrations to the database
// at the specified database URL.
func tryToMigrate(dbURL string) error {
	source := migrationsSource()
	m, err := migrate.New(source, dbURL)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("m.Up: %w", err)
	}
	return nil
}

// migrationsSource returns a uri pointing to the migrations directory.
// it panics on failure.
func migrationsSource() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to determine relative path")
	}

	return "file://" + filepath.Clean(filepath.Join(filepath.Dir(filename), filepath.FromSlash("migrations")))
}
