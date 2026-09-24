package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithPanels(builder *dashboard.DashboardBuilder) {
	// RPS
	builder.WithPanel(WithPanelSize(g.ServerRPS(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.ClientRPS(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.QueueIngoingMPS(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.QueueOutgoingMPS(), PanelSizeQuarterHigh))
	// Latency
	builder.WithPanel(WithPanelSize(g.ServerLatency(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.ClientLatency(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.QueueIngoingLatency(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.QueueOutgoingLatency(), PanelSizeQuarterHigh))
	// Error rate
	// Gorutines
}
