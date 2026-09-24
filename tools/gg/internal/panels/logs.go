package panels

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/variants"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
)

func (g Generator) Logs() *logs.PanelBuilder {
	return logs.
		NewPanelBuilder().
		Title("Logs").
		Targets([]cog.Builder[variants.Dataquery]{
			loki. // Примечание по сигнатуре частично совпадает, т.ч. используем его.
				NewDataqueryBuilder().
				Expr(g.core.LogInstanceFilter()),
		}).
		SortOrder(common.LogsSortOrderDescending).
		EnableLogDetails(true).
		Datasource(g.core.LogsDatasource())
}
