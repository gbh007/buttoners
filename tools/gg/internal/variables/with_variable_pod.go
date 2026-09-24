package variables

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithVariablePod(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(
		dashboard.
			NewQueryVariableBuilder("pod").
			Query(dashboard.StringOrMap{
				String: new(`label_values(process_cpu_seconds_total{service=~"$service"}, pod)`),
			}).
			Datasource(g.core.MetricDatasource()).
			IncludeAll(true).
			AllValue(".+").
			Multi(true).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged),
	)
}
