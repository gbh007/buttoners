package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueIngoingErrorRate() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Входящий error rate очередей").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_queue_reader_handle_seconds",
					[]string{"server_addr"},
					`status="err"`,
				)).
				LegendFormat("queue reader => {{server_addr}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.PercentUnit).
		Datasource(g.core.MetricsDatasource())
}
