package panels

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ServerRPS() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Входящий RPS").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_grpc_server_handle_seconds_count{%s}[$__rate_interval])) by (server_addr)",
					g.core.InstanceFilter(),
				)).
				LegendFormat("grpc server => {{server_addr}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_http_server_handle_seconds_count{%s}[$__rate_interval])) by (server_addr)",
					g.core.InstanceFilter(),
				)).
				LegendFormat("http server => {{server_addr}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.RequestsPerSecond).
		Datasource(g.core.MetricDatasource())
}
