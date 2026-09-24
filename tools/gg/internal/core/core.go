package core

import (
	"fmt"
	"slices"
	"strings"

	"github.com/grafana/grafana-foundation-sdk/go/common"
)

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

func (c Core) LatencyWithInstanceFilter(metricName string, additionalLables []string) string {
	return fmt.Sprintf(
		`histogram_quantile($quantile, sum(rate(%s_bucket{%s}[$__rate_interval])) by (%s))`,
		metricName,
		c.InstanceFilter(),
		strings.Join(append(slices.Clip(additionalLables), "le"), ", "),
	)
}

func (c Core) RPSWithInstanceFilterFromHistogram(metricName string, additionalLables []string) string {
	return fmt.Sprintf(
		`sum(rate(%s_count{%s}[$__rate_interval])) by (%s)`,
		metricName,
		c.InstanceFilter(),
		strings.Join(additionalLables, ", "),
	)
}
