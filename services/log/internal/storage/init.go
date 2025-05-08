package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/gbh007/buttoners/services/log/internal/storage/migration"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"

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

	goose.SetBaseFS(migration.Migrations)

	err = goose.SetDialect(string(goose.DialectClickHouse))
	if err != nil {
		return nil, fmt.Errorf("%w: set dialect: %w", errDatabase, err)
	}

	err = goose.UpContext(
		ctx, db.DB, ".",
		goose.WithNoColor(true),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: up migrations: %w", errDatabase, err)
	}

	return &Database{
		db: db,
	}, nil
}
