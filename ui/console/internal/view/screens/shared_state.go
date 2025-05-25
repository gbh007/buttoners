package screens

import (
	"context"

	"github.com/gbh007/buttoners/ui/console/internal/storage"
)

type SharedState struct {
	Ctx     context.Context
	Storage *storage.Storage
}
