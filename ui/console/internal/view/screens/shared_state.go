package screens

import (
	"context"

	"github.com/gbh007/buttoners/core/clients/gateclient"
)

type SharedState struct {
	Ctx        context.Context
	GateClient *gateclient.Client
	GateToken  string
}
