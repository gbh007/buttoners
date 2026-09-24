package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func (g Generator) QueueOutgoingErrorRate() *timeseries.PanelBuilder {
	targets := []cog.Builder[variants.Dataquery]{}

	if g.core.HasModule(core.ModuleQueueWriter) {
		targets = append(
			targets,
			prometheus.
				NewDataqueryBuilder().
				Expr(g.core.ErrorRateWithInstanceFilterFromHistogram(
					"buttoners_queue_writer_handle_seconds",
					[]string{"target_host"},
					`status="err"`,
				)).
				LegendFormat("queue writer => {{target_host}}"),
		)
	}

	if len(targets) == 0 {
		panic("unexpected behavior")
	}

	return timeseries.
		NewPanelBuilder().
		Title("Исходящий error rate очередей").
		Targets(targets).
		Legend(g.core.SimpleLegend()).
		Unit(units.PercentUnit).
		Datasource(g.core.MetricsDatasource())
}
