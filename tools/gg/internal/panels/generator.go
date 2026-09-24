package panels

import "github.com/gbh007/buttoners/tools/gg/internal/core"

type Generator struct {
	core core.Core
}

func New(core core.Core) Generator {
	return Generator{
		core: core,
	}
}
