package panels

import (
	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
)

func (g Generator) Logs() *logs.PanelBuilder {
	targets := []cog.Builder[variants.Dataquery]{}

	if g.core.HasModule(core.ModuleLogs) {
		targets = append(
			targets,
			loki. // Примечание по сигнатуре частично совпадает, т.ч. используем его.
				NewDataqueryBuilder().
				Expr(g.core.LogInstanceFilter()),
		)
	}

	if len(targets) == 0 {
		panic("unexpected behavior")
	}

	return logs.
		NewPanelBuilder().
		Title("Logs").
		Targets(targets).
		SortOrder(common.LogsSortOrderDescending).
		EnableLogDetails(true).
		Datasource(g.core.LogsDatasource())
}
