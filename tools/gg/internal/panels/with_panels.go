package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func (g Generator) WithPanels(builder *dashboard.DashboardBuilder) {
	// RPS

	if g.core.HasAnyModule(core.ModuleGRPCServer, core.ModuleHTTPServer) {
		builder.WithPanel(WithPanelSize(g.ServerRPS(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleGRPCClient, core.ModuleHTTPClient, core.ModuleRedis) {
		builder.WithPanel(WithPanelSize(g.ClientRPS(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleQueueReader) {
		builder.WithPanel(WithPanelSize(g.QueueIngoingMPS(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleQueueWriter) {
		builder.WithPanel(WithPanelSize(g.QueueOutgoingMPS(), PanelSizeQuarterHigh))
	}

	// Latency

	if g.core.HasAnyModule(core.ModuleGRPCServer, core.ModuleHTTPServer) {
		builder.WithPanel(WithPanelSize(g.ServerLatency(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleGRPCClient, core.ModuleHTTPClient, core.ModuleRedis) {
		builder.WithPanel(WithPanelSize(g.ClientLatency(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleQueueReader) {
		builder.WithPanel(WithPanelSize(g.QueueIngoingLatency(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleQueueWriter) {
		builder.WithPanel(WithPanelSize(g.QueueOutgoingLatency(), PanelSizeQuarterHigh))
	}

	// Error rate

	if g.core.HasAnyModule(core.ModuleGRPCServer, core.ModuleHTTPServer) {
		builder.WithPanel(WithPanelSize(g.ServerErrorRate(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleGRPCClient, core.ModuleHTTPClient, core.ModuleRedis) {
		builder.WithPanel(WithPanelSize(g.ClientErrorRate(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleQueueReader) {
		builder.WithPanel(WithPanelSize(g.QueueIngoingErrorRate(), PanelSizeQuarterHigh))
	}

	if g.core.HasAnyModule(core.ModuleQueueWriter) {
		builder.WithPanel(WithPanelSize(g.QueueOutgoingErrorRate(), PanelSizeQuarterHigh))
	}

	// Logs

	if g.core.HasAnyModule(core.ModuleLogs) {
		builder.WithPanel(WithPanelSize(g.Logs(), PanelSizeFull))
	}

	// Traces

	if g.core.HasAnyModule(core.ModuleTraces) {
		builder.WithPanel(WithPanelSize(g.Traces(), PanelSizeFull))
	}
}
