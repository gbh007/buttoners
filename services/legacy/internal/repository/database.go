package repository

import (
	"github.com/gbh007/buttoners/services/legacy/internal/domain"
	"fmt"
	"log/slog"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
}

func New(lg *slog.Logger, dbType, dbDNS string) (*Repository, error) {
	var dialector gorm.Dialector

	switch dbType {
	case "sqlite":
		dialector = sqlite.Open(dbDNS)
	case "postgres":
		dialector = postgres.Open(dbDNS)
	case "mysql":
		dialector = mysql.Open(dbDNS)
	default:
		return nil, fmt.Errorf("unknown db type: %s", dbType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{Logger: logger.New(
		slog.NewLogLogger(lg.Handler(), slog.LevelDebug),
		logger.Config{LogLevel: logger.Info},
	)})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Button{},
	)
	if err != nil {
		return nil, fmt.Errorf("gorm automigrate: %w", err)
	}

	return &Repository{
		db: db,
	}, nil
}
