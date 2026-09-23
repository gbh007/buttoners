package internal

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

// Базовая метрика для определения контейнера
// process_cpu_seconds_total
// process_cpu_seconds_total {instance="172.20.0.17:8082",job="prometheus.scrape.container_metrics",pod="buttoners-auth-1",service="auth"}

type Generator struct {
	uid      string
	services []string
}

func New(uid string, services []string) *Generator {
	plugins.RegisterDefaultPlugins()

	return &Generator{
		uid:      uid,
		services: services,
	}
}

func (g Generator) Build() (dashboard.Dashboard, error) {
	builder := dashboard.
		NewDashboardBuilder("Buttoners").
		Uid(g.uid).
		Timezone("Asia/Krasnoyarsk").
		Time("now-4h", "now").
		WeekStart("monday").
		Refresh("1m").
		Tooltip(dashboard.DashboardCursorSyncCrosshair)

	g.WithPanels(builder)
	g.WithVariables(builder)

	d, err := builder.Build()
	if err != nil {
		return dashboard.Dashboard{}, fmt.Errorf("build: %w", err)
	}

	return d, nil
}

func (g Generator) WithVariables(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
	builder.WithVariable(
		dashboard.
			NewDatasourceVariableBuilder("metrics").
			Type("prometheus"),
	)

	builder.WithVariable(
		dashboard.
			NewQueryVariableBuilder("service").
			Query(dashboard.StringOrMap{
				String: new(`label_values(process_cpu_seconds_total, service)`),
			}).
			Datasource(metricDS).
			IncludeAll(true).
			AllValue(".+").
			Multi(true).
			Current(dashboard.VariableOption{
				Selected: new(true),
				Text: dashboard.StringOrArrayOfString{
					ArrayOfString: g.services,
				},
				Value: dashboard.StringOrArrayOfString{
					ArrayOfString: g.services,
				},
			}).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged),
	)

	builder.WithVariable(
		dashboard.
			NewQueryVariableBuilder("pod").
			Query(dashboard.StringOrMap{
				String: new(`label_values(process_cpu_seconds_total{service=~"$service"}, pod)`),
			}).
			Datasource(metricDS).
			IncludeAll(true).
			AllValue(".+").
			Multi(true).
			Current(dashboard.VariableOption{
				Selected: new(true),
				Text: dashboard.StringOrArrayOfString{
					ArrayOfString: g.services,
				},
				Value: dashboard.StringOrArrayOfString{
					ArrayOfString: g.services,
				},
			}).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged),
	)

	return builder
}

func (g Generator) WithPanels(builder *dashboard.DashboardBuilder) *dashboard.DashboardBuilder {
	builder.WithPanel(timeseries.
		NewPanelBuilder().
		Title("RPS").
		Targets([]cog.Builder[variants.Dataquery]{
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_grpc_client_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					serviceFilter,
				)).
				LegendFormat("grpc client => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_grpc_server_handle_seconds_count{%s}[$__rate_interval])) by (server_addr)",
					serviceFilter,
				)).
				LegendFormat("grpc server => {{server_addr}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_http_client_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					serviceFilter,
				)).
				LegendFormat("http client => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_http_server_handle_seconds_count{%s}[$__rate_interval])) by (server_addr)",
					serviceFilter,
				)).
				LegendFormat("http server => {{server_addr}}"),

			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_queue_writer_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					serviceFilter,
				)).
				LegendFormat("queue writer => {{target_host}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_queue_reader_handle_seconds_count{%s}[$__rate_interval])) by (server_addr)",
					serviceFilter,
				)).
				LegendFormat("queue reader => {{server_addr}}"),
			prometheus.
				NewDataqueryBuilder().
				Expr(fmt.Sprintf(
					"sum(rate(buttoners_redis_handle_seconds_count{%s}[$__rate_interval])) by (target_host)",
					serviceFilter,
				)).
				LegendFormat("redis => {{target_host}}"),
		}).
		Legend(simpleLegend()).
		Unit(units.RequestsPerSecond).
		Datasource(metricDS))

	return builder
}
