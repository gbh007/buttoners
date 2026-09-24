package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueIngoingLatency() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Входящие задержки очередей").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.LatencyWithInstanceFilter(
					"buttoners_queue_reader_handle_seconds",
					[]string{"server_addr"},
				)).
				LegendFormat("queue reader => {{server_addr}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.Seconds).
		Datasource(g.core.MetricsDatasource())
}
