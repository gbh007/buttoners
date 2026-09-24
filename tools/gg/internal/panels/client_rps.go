package panels

import (
	"fmt"

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
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_grpc_client_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					g.core.InstanceFilter(),
				)).
				LegendFormat("grpc client => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_http_client_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					g.core.InstanceFilter(),
				)).
				LegendFormat("http client => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_redis_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					g.core.InstanceFilter(),
				)).
				LegendFormat("redis => {{target_host}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.RequestsPerSecond).
		Datasource(g.core.MetricDatasource())
}
