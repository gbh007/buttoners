package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueOutgoingErrorRate() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Исходящий error rate очередей").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_queue_writer_handle_seconds",
					[]string{"target_host"},
					`status="err"`,
				)).
				LegendFormat("queue writer => {{target_host}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.PercentUnit).
		Datasource(g.core.MetricDatasource())
}
