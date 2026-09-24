package panels

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueIngoingMPS() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Входящий MPS").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_queue_reader_handle_seconds_count{%s}[$__rate_interval])) by (server_addr)",
					g.core.InstanceFilter(),
				)).
				LegendFormat("queue reader => {{server_addr}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.MessagesPerSecond).
		Datasource(g.core.MetricDatasource())
}
