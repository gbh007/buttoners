package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gbh007/buttoners/services/worker/internal/storage/migration"
	"github.com/jmoiron/sqlx"
	_ "github.com/mailru/go-clickhouse/v2" // imp[ort sql driver
	"github.com/pressly/goose/v3"
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

	goose.SetBaseFS(migration.Migrations)

	err = goose.SetDialect(string(goose.DialectClickHouse))
	if err != nil {
		return nil, fmt.Errorf("%w: set dialect: %w", errDatabase, err)
	}

	err = goose.UpContext(
		ctx, conn, ".",
		goose.WithNoColor(true),
		goose.WithAllowMissing(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: up migrations: %w", errDatabase, err)
	}

	return &Database{
		db: db,
	}, nil
}
