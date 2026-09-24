package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/table"
	"github.com/grafana/grafana-foundation-sdk/go/tempo"
)

func (g Generator) Traces() *table.PanelBuilder {
	return table.
		NewPanelBuilder().
		Title("Traces").
		Targets([]cog.Builder[variants.Dataquery]{
			tempo.
				NewDataqueryBuilder().
				QueryType("traceql").
				Query(g.core.TraceInstanceFilter()),
		}).
		Datasource(g.core.TracesDatasource())
}
