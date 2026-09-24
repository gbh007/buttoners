package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/table"
	"github.com/grafana/grafana-foundation-sdk/go/tempo"
)

func (g Generator) Traces() *table.PanelBuilder {
	targets := []cog.Builder[variants.Dataquery]{}

	if g.core.HasModule(core.ModuleTraces) {
		targets = append(
			targets,
			tempo.
				NewDataqueryBuilder().
				QueryType("traceql").
				Query(g.core.TraceInstanceFilter()),
		)
	}

	if len(targets) == 0 {
		panic("unexpected behavior")
	}

	return table.
		NewPanelBuilder().
		Title("Traces").
		Targets(targets).
		Datasource(g.core.TracesDatasource())
}
