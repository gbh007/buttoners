package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithPanels(builder *dashboard.DashboardBuilder) {
	builder.WithPanel(g.ServerRPS())
	builder.WithPanel(g.ClientRPS())
	builder.WithPanel(g.QueueIngoingMPS())
	builder.WithPanel(g.QueueOutgoingMPS())
}
