package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueOutgoingLatency() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Исходящие задержки очередей").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.LatencyWithInstanceFilter(
					"buttoners_queue_writer_handle_seconds",
					[]string{"target_host"},
				)).
				LegendFormat("queue writer => {{target_host}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.MessagesPerSecond).
		Datasource(g.core.MetricsDatasource())
}
