package variables

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithVariableQuantile(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(
		dashboard.
			NewCustomVariableBuilder("quantile").
			Current(dashboard.VariableOption{
				Selected: new(true),
				Text: dashboard.StringOrArrayOfString{
					String: new(".80"),
				},
				Value: dashboard.StringOrArrayOfString{
					String: new(".80"),
				},
			}).
			Values(dashboard.StringOrMap{
				String: new(".50,.80,.95,.99"),
			}),
	)
}
