package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ServerRPS() *timeseries.PanelBuilder {
	targets := []cog.Builder[variants.Dataquery]{}

	if g.core.HasModule(core.ModuleGRPCServer) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_grpc_server_handle_seconds",
					[]string{"server_addr"},
				)).
				LegendFormat("grpc server => {{server_addr}}"),
		)
	}

	if g.core.HasModule(core.ModuleHTTPServer) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.RPSWithInstanceFilterFromHistogram(
					"buttoners_http_server_handle_seconds",
					[]string{"server_addr"},
				)).
				LegendFormat("http server => {{server_addr}}"),
		)
	}

	if len(targets) == 0 {
		panic("unexpected behavior")
	}

	return timeseries.
		NewPanelBuilder().
		Title("Входящий RPS").
		Targets(targets).
		Legend(g.core.SimpleLegend()).
		Unit(units.RequestsPerSecond).
		Datasource(g.core.MetricsDatasource())
}
