package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/gbh007/buttoners/services/auth/internal/storage/migration"
	_ "github.com/go-sql-driver/mysql" // import sql driver
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

var errDatabase = errors.New("auth database")

type Database struct {
	db *sqlx.DB
}

func New(ctx context.Context, username, password, dbHostWithPort, databaseName string) (*Database, error) {
	cs := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true",
		username, password, dbHostWithPort, databaseName,
	)

	db, err := otelsqlx.Open("mysql", cs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabase, err)
	}

	goose.SetBaseFS(migration.Migrations)

	err = goose.SetDialect(string(goose.DialectMySQL))
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
