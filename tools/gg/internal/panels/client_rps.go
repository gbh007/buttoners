package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ClientRPS() *timeseries.PanelBuilder {
	targets := []cog.Builder[variants.Dataquery]{}

	if g.core.HasModule(core.ModuleGRPCClient) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_grpc_client_handle_seconds",
					[]string{"target_host"},
				)).
				LegendFormat("grpc client => {{target_host}}"),
		)
	}

	if g.core.HasModule(core.ModuleHTTPClient) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_http_client_handle_seconds",
					[]string{"target_host"},
				)).
				LegendFormat("http client => {{target_host}}"),
		)
	}

	if g.core.HasModule(core.ModuleRedis) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_redis_handle_seconds",
					[]string{"target_host"},
				)).
				LegendFormat("redis => {{target_host}}"),
		)
	}

	if len(targets) == 0 {
		panic("unexpected behavior")
	}

	return timeseries.
		NewPanelBuilder().
		Title("Исходящий RPS").
		Targets(targets).
		Legend(g.core.SimpleLegend()).
		Unit(units.RequestsPerSecond).
		Datasource(g.core.MetricsDatasource())
}
