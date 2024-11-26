package migration

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func MigrateDatabase(connectionString string) error {
	m, err := migrate.New(
		"file:..//db//migrations",
		connectionString,
	)
	if err != nil {
		log.Println("Migration failed with error:", err)
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println("Up migration failed with error:", err)
		return err
	}
	return nil
}
