package variables

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithVariableTraces(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(
		dashboard.
			NewDatasourceVariableBuilder("traces").
			Hide(dashboard.VariableHideHideVariable).
			Type("tempo"),
	)
}
