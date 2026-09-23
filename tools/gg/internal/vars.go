package internal

import "github.com/grafana/grafana-foundation-sdk/go/common"

const serviceFilter = `service=~"$service", pod=~"$pod"`

var metricDS = common.DataSourceRef{
	Type: new("prometheus"),
	Uid:  new("${metrics}"),
}

func simpleLegend() *common.VizLegendOptionsBuilder {
	return common.
		NewVizLegendOptionsBuilder().
		DisplayMode(common.LegendDisplayModeTable).
		Placement(common.LegendPlacementBottom).
		Calcs([]string{"mean", "lastNotNull"}).
		SortBy("Mean").
		SortDesc(true).
		ShowLegend(true)
}
