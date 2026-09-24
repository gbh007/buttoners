package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueIngoingLatency() *timeseries.PanelBuilder {
	targets := []cog.Builder[variants.Dataquery]{}

	if g.core.HasModule(core.ModuleQueueReader) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.LatencyWithInstanceFilter(
					"buttoners_queue_reader_handle_seconds",
					[]string{"server_addr"},
				)).
				LegendFormat("queue reader => {{server_addr}}"),
		)
	}

	if len(targets) == 0 {
		panic("unexpected behavior")
	}

	return timeseries.
		NewPanelBuilder().
		Title("Входящие задержки очередей $quantile").
		Targets(targets).
		Legend(g.core.SimpleLegend()).
		Unit(units.Seconds).
		Datasource(g.core.MetricsDatasource())
}
