package core

import (
	"fmt"
	"slices"
	"strings"
)

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

func (c Core) ErrorRateWithInstanceFilterFromHistogram(metricName string, additionalLables []string, errorFilter string) string {
	return fmt.Sprintf(
		`sum(rate(%s_count{%s, %s}[$__rate_interval])) by (%s) or vector(0)/ sum(rate(%s_count{%s}[$__rate_interval])) by (%s)`,
		metricName,
		c.InstanceFilter(),
		errorFilter,
		strings.Join(additionalLables, ", "),
		metricName,
		c.InstanceFilter(),
		strings.Join(additionalLables, ", "),
	)
}
