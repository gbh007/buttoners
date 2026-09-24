package variables

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithVariables(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(
		dashboard.
			NewDatasourceVariableBuilder("metrics").
			Type("prometheus"),
	)

	g.WithVariableService(builder)
	g.WithVariablePod(builder)
	g.WithVariableQuantile(builder)
}
