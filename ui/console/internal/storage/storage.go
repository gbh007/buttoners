package storage

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

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
	// FIXME: удалить после отладки
	time.Sleep(time.Millisecond * time.Duration(rand.IntN(1000)+500))
	if rand.IntN(2) == 0 {
		return fmt.Errorf("ooooops")
	}

	s.data.Connection = data
	return nil
}

func (s *Storage) GetConnectionData(ctx context.Context) (model.Connection, error) {
	// FIXME: удалить после отладки
	time.Sleep(time.Millisecond * time.Duration(rand.IntN(1000)+500))
	if rand.IntN(2) == 0 {
		return model.Connection{}, fmt.Errorf("ooooops")
	}

	return s.data.Connection, nil
}
