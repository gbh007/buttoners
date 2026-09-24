package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ServerErrorRate() *timeseries.PanelBuilder {
	targets := []cog.Builder[variants.Dataquery]{}

	if g.core.HasModule(core.ModuleGRPCServer) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_grpc_server_handle_seconds",
					[]string{"server_addr"},
					`status!="OK"`,
				)).
				LegendFormat("grpc server => {{server_addr}}"),
		)
	}

	if g.core.HasModule(core.ModuleHTTPServer) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_http_server_handle_seconds",
					[]string{"server_addr"},
					`status=~"4\\d{2}"`,
				)).
				LegendFormat("http server 4xx => {{server_addr}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_http_server_handle_seconds",
					[]string{"server_addr"},
					`status=~"5\\d{2}"`,
				)).
				LegendFormat("http server 5xx => {{server_addr}}"),
		)
	}

	if len(targets) == 0 {
		panic("unexpected behavior")
	}

	return timeseries.
		NewPanelBuilder().
		Title("Входящий error rate").
		Targets(targets).
		Legend(g.core.SimpleLegend()).
		Unit(units.PercentUnit).
		Datasource(g.core.MetricsDatasource())
}
