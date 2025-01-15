package database

import (
	"OnlineMusic/config"
	migration "OnlineMusic/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type Database struct {
	Connection *pgx.Conn
}

func New(ctx context.Context, cfg config.Database) (*Database, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
	slog.Info("Connecting to database", slog.String("connection string", connectionString))
	connection, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	if err := migration.MigrateDatabase(connectionString); err != nil {
		slog.With("error", err).Error("Failed to apply migrations")
		return nil, err
	}
	return &Database{Connection: connection}, nil
}

func (db *Database) Close(ctx context.Context) error {
	return db.Connection.Close(ctx)
}
