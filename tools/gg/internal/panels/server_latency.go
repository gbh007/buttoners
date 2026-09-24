package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ServerLatency() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Входящие задержки").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.LatencyWithInstanceFilter(
					"buttoners_grpc_server_handle_seconds",
					[]string{"server_addr"},
				)).
				LegendFormat("grpc server => {{server_addr}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.LatencyWithInstanceFilter(
					"buttoners_http_server_handle_seconds",
					[]string{"server_addr"},
				)).
				LegendFormat("http server => {{server_addr}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.RequestsPerSecond).
		Datasource(g.core.MetricDatasource())
}
