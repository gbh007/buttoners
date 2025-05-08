package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/gbh007/buttoners/services/log/internal/storage/migration"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"

	"github.com/golang-migrate/migrate/v4"
	cii "github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mailru/go-clickhouse/v2" // import sql driver
)

var errDatabase = errors.New("log database")

type Database struct {
	db *sqlx.DB
}

func Init(ctx context.Context, username, password, dbHostWithPort, databaseName string) (*Database, error) {
	cs := fmt.Sprintf("http://%s:%s@%s/%s", username, password, dbHostWithPort, databaseName)

	db, err := otelsqlx.Open("chhttp", cs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabase, err)
	}

	sourceInstance, err := iofs.New(migration.Migrations, ".")
	if err != nil {
		return nil, fmt.Errorf("%w: open source: %w", errDatabase, err)
	}

	dbInstance, err := cii.WithInstance(db.DB, &cii.Config{
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
	if err != nil {
		return nil, fmt.Errorf("%w: up migrations: %w", errDatabase, err)
	}

	return &Database{
		db: db,
	}, nil
}
