package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ClientRPS() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Исходящий RPS").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_grpc_client_handle_seconds",
					[]string{"target_host"},
				)).
				LegendFormat("grpc client => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_http_client_handle_seconds",
					[]string{"target_host"},
				)).
				LegendFormat("http client => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_redis_handle_seconds",
					[]string{"target_host"},
				)).
				LegendFormat("redis => {{target_host}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.RequestsPerSecond).
		Datasource(g.core.MetricsDatasource())
}
