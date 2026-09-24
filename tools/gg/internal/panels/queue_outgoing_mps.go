package panels

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueOutgoingMPS() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Исходящий MPS").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_queue_writer_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					g.core.InstanceFilter(),
				)).
				LegendFormat("queue writer => {{target_host}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.MessagesPerSecond).
		Datasource(g.core.MetricDatasource())
}
