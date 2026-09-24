package variables

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

// FIXME: инвертировать переменные на подобии панелей
func (g Generator) WithVariables(builder *dashboard.DashboardBuilder) {
	g.WithVariableMetrics(builder)
	g.WithVariableLogs(builder)
	g.WithVariableService(builder)
	g.WithVariablePod(builder)
	g.WithVariableQuantile(builder)
	g.WithVariableLogs(builder)
}
