package variables

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithVariableLogs(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(
		dashboard.
			NewDatasourceVariableBuilder("logs").
			Hide(dashboard.VariableHideHideVariable).
			Type("victoriametrics-logs-datasource"),
	)
}
