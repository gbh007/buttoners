package variables

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithVariableMetrics(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(
		dashboard.
			NewDatasourceVariableBuilder("metrics").
			Hide(dashboard.VariableHideHideVariable).
			Type("prometheus"),
	)
}
