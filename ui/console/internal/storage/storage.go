package storage

import (
	"context"

	"github.com/gbh007/buttoners/ui/console/internal/model"
)

type _ConfigData struct {
	Connection model.Connection `json:"connection"`
}

type Storage struct {
	data _ConfigData
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Load(ctx context.Context) error {
	return nil
}

func (s *Storage) Save(ctx context.Context) error {
	return nil
}

func (s *Storage) SetConnectionData(ctx context.Context, data model.Connection) error {
	s.data.Connection = data
	return nil
}

func (s *Storage) GetConnectionData(ctx context.Context) (model.Connection, error) {
	return s.data.Connection, nil
}
