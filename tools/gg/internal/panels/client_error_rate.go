package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) ClientErrorRate() *timeseries.PanelBuilder {
	return timeseries.
		NewPanelBuilder().
		Title("Исходящий error rate").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_grpc_client_handle_seconds",
					[]string{"target_host"},
					`status!="OK"`,
				)).
				LegendFormat("grpc client => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_http_client_handle_seconds",
					[]string{"target_host"},
					`status=~"4\\d{2}"`,
				)).
				LegendFormat("http client 4xx => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_http_client_handle_seconds",
					[]string{"target_host"},
					`status=~"5\\d{2}"`,
				)).
				LegendFormat("http client 5xx => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_redis_handle_seconds",
					[]string{"target_host"},
					`status="err"`,
				)).
				LegendFormat("redis => {{target_host}}"),
		}).
		Legend(g.core.SimpleLegend()).
		Unit(units.PercentUnit).
		Datasource(g.core.MetricDatasource())
}
