package variables

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithVariableService(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(
		dashboard.
			NewQueryVariableBuilder("service").
			Query(dashboard.StringOrMap{
				String: new(fmt.Sprintf(
					`label_values(process_cpu_seconds_total{service=~"%s"}, service)`,
					g.core.DirtyServiceFilter,
				)),
			}).
			Datasource(g.core.MetricDatasource()).
			IncludeAll(true).
			AllValue(".+").
			Multi(true).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged),
	)
}
