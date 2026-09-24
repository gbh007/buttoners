package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ServerErrorRate() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Входящий error rate").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_grpc_server_handle_seconds",
					[]string{"server_addr"},
					`status!="OK"`,
				)).
				LegendFormat("grpc server => {{server_addr}}"),
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
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.PercentUnit).
		Datasource(g.core.MetricDatasource())
}
