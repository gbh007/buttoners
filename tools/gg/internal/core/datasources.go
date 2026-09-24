package core

import "github.com/grafana/grafana-foundation-sdk/go/common"

func (Core) MetricsDatasource() common.DataSourceRef {
	return common.DataSourceRef{
		Type: new("prometheus"),
		Uid:  new("${metrics}"),
	}
}

func (Core) LogsDatasource() common.DataSourceRef {
	return common.DataSourceRef{
		Type: new("victoriametrics-logs-datasource"),
		Uid:  new("${logs}"),
	}
}

func (Core) TracesDatasource() common.DataSourceRef {
	return common.DataSourceRef{
		Type: new("tempo"),
		Uid:  new("${traces}"),
	}
}
