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
	builder.WithPanel(WithPanelSize(g.ServerErrorRate(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.ClientErrorRate(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.QueueIngoingErrorRate(), PanelSizeQuarterHigh))
	builder.WithPanel(WithPanelSize(g.QueueOutgoingErrorRate(), PanelSizeQuarterHigh))
	// Logs
	builder.WithPanel(WithPanelSize(g.Logs(), PanelSizeFull))
	// Traces
	builder.WithPanel(WithPanelSize(g.Traces(), PanelSizeFull))
}
