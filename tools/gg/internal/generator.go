package internal

import (
	"fmt"

	"github.com/gbh007/buttoners/tools/gg/internal/core"
	"github.com/gbh007/buttoners/tools/gg/internal/panels"
	"github.com/gbh007/buttoners/tools/gg/internal/variables"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

// Базовая метрика для определения контейнера
// process_cpu_seconds_total
// process_cpu_seconds_total {instance="172.20.0.17:8082",job="prometheus.scrape.container_metrics",pod="buttoners-auth-1",service="auth"}

type Generator struct {
	uid  string
	name string

	core core.Core
}

func New(uid string, name string, dirtyServiceFilter string, modules []string) Generator {
	plugins.RegisterDefaultPlugins()

	return Generator{
		uid:  uid,
		name: name,
		core: core.Core{
			DirtyServiceFilter: dirtyServiceFilter,
			EnabledModules:     modules,
		},
	}
}

func (g Generator) Build() (dashboard.Dashboard, error) {
	builder := dashboard.
		NewDashboardBuilder(g.name).
		Uid(g.uid).
		Timezone("Asia/Krasnoyarsk").
		Time("now-4h", "now").
		WeekStart("monday").
		Refresh("1m").
		Tags([]string{"buttoners"}).
		Links([]cog.Builder[dashboard.DashboardLink]{
			dashboard.NewDashboardLinkBuilder("Buttoners boards").
				AsDropdown(true).
				KeepTime(true).
				Type(dashboard.DashboardLinkTypeDashboards).
				Tags([]string{"buttoners"}),
			dashboard.NewDashboardLinkBuilder("Github").
				TargetBlank(true).
				Type(dashboard.DashboardLinkTypeLink).
				Url("https://github.com/gbh007/buttoners"),
		}).
		Tooltip(dashboard.DashboardCursorSyncCrosshair)

	variables.New(g.core).WithVariables(builder)
	panels.New(g.core).WithPanels(builder)

	d, err := builder.Build()
	if err != nil {
		return dashboard.Dashboard{}, fmt.Errorf("build: %w", err)
	}

	return d, nil
}
