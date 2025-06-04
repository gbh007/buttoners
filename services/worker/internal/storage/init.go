package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gbh007/buttoners/services/worker/internal/storage/migration"
	"github.com/golang-migrate/migrate/v4"
	cii "github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mailru/go-clickhouse/v2" // imp[ort sql driver // FIXME: заменить на другую либу
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

var errDatabase = errors.New("task database")

type Database struct {
	db *sqlx.DB
}

func Init(ctx context.Context, username, password, dbHostWithPort, databaseName string) (*Database, error) {
	cs := fmt.Sprintf("http://%s:%s@%s/%s", username, password, dbHostWithPort, databaseName)

	db, err := otelsqlx.Open("chhttp", cs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabase, err)
	}

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Protocol: clickhouse.HTTP,
		Addr:     []string{dbHostWithPort},
		Auth: clickhouse.Auth{
			Database: databaseName,
			Username: username,
			Password: password,
		},
	})

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w: ping: %w", errDatabase, err)
	}

	sourceInstance, err := iofs.New(migration.Migrations, ".")
	if err != nil {
		return nil, fmt.Errorf("%w: open source: %w", errDatabase, err)
	}

	dbInstance, err := cii.WithInstance(conn, &cii.Config{
		DatabaseName:          databaseName,
		MigrationsTable:       "my_migrations",
		MigrationsTableEngine: "MergeTree",
	})
	if err != nil {
		return nil, fmt.Errorf("%w: open source: %w", errDatabase, err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceInstance,
		"clickhouse",
		dbInstance,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: new migrate: %w", errDatabase, err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%w: up migrations: %w", errDatabase, err)
	}

	return &Database{
		db: db,
	}, nil
}
