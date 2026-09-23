package main

import (
	"encoding/json"
	"flag"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"

	"github.com/BurntSushi/toml"
	generator "github.com/gbh007/buttoners/tools/gg/internal"
)

type Config struct {
	Grafana struct {
		Addr   string `toml:"addr"`
		Token  string `toml:"token"`
		Folder string `toml:"folder"`
		UID    string `toml:"uid"`
	} `toml:"grafana"`
	To       string   `toml:"to"`
	Services []string `toml:"services"`
}

func main() {
	configPath := flag.String("config", "gg-config.toml", "path to config")

	flag.Parse()

	var cfg Config

	f, err := os.Open(*configPath)
	if err != nil {
		panic(err)
	}

	_, err = toml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Grafana.UID == "" {
		panic("empty uid")
	}

	g := generator.New(cfg.Grafana.UID, cfg.Services)

	dashboardModel, err := g.Build()
	if err != nil {
		panic(err)
	}

	if cfg.Grafana.Addr != "" {
		u, err := url.Parse(cfg.Grafana.Addr)
		if err != nil {
			panic(err)
		}

		transportCfg := &goapi.TransportConfig{
			Host:     u.Host,
			BasePath: u.Path,
			Schemes:  []string{u.Scheme},
			APIKey:   cfg.Grafana.Token,
		}

		client := goapi.NewHTTPClientWithConfig(strfmt.Default, transportCfg)

		response, err := client.Dashboards.PostDashboard(&models.SaveDashboardCommand{
			FolderUID: cfg.Grafana.Folder,
			Dashboard: dashboardModel,
			Overwrite: true,
		})
		if err != nil {
			panic(err)
		}

		if *response.Payload.Status != "success" {
			panic(*response.Payload.Status)
		}
	}

	if cfg.To != "" {
		out, err := os.Create(cfg.To)
		if err != nil {
			panic(err)
		}

		enc := json.NewEncoder(out)
		enc.SetIndent("", "   ")

		err = enc.Encode(dashboardModel)
		if err != nil {
			panic(err)
		}

		err = out.Close()
		if err != nil {
			panic(err)
		}
	}
}
