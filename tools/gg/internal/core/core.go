package core

import "github.com/grafana/grafana-foundation-sdk/go/common"

type Core struct {
	DirtyServiceFilter string
}

func (Core) SimpleLegend() *common.VizLegendOptionsBuilder {
	return common.
		NewVizLegendOptionsBuilder().
		DisplayMode(common.LegendDisplayModeTable).
		Placement(common.LegendPlacementBottom).
		Calcs([]string{"mean", "lastNotNull"}).
		SortBy("Mean").
		SortDesc(true).
		ShowLegend(true)
}

func (Core) MetricDatasource() common.DataSourceRef {
	return common.DataSourceRef{
		Type: new("prometheus"),
		Uid:  new("${metrics}"),
	}
}

func (Core) InstanceFilter() string {
	return `service=~"$service", pod=~"$pod"`
}
