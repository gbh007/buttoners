package controller

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/services/legacy/internal/repository"
	"github.com/gbh007/buttoners/services/legacy/internal/service/button"
	"github.com/gbh007/buttoners/services/legacy/internal/service/user"
	"github.com/go-playground/validator/v10"

	cMetrics "github.com/gbh007/buttoners/core/metrics"
	"github.com/labstack/echo/v4"

	_ "embed"
)

type Controller struct {
	addr  string
	debug bool

	logger *slog.Logger

	buttonService *button.Service
	userSevice    *user.Service
}

type Config struct {
	APIAddr string `envconfig:"SELF_ADDR" default:":8080"`
	Debug   bool   `envconfig:"SELF_DEBUG"`
	DBType  string `envconfig:"DB_TYPE"`
	DBDNS   string `envconfig:"DB_DNS"`

	AuthAddr  string `envconfig:"AUTH_SERVICE_ADDR"`
	AuthToken string `envconfig:"AUTH_SERVICE_TOKEN"`

	JaegerURL      string `envconfig:"JAEGER_URL" default:"http://jaeger:14268/api/traces"`
	PrometheusAddr string `envconfig:"PROMETHEUS_ADDR" default:"pushgateway:9091"`
}

func New(logger *slog.Logger, cfg Config) (*Controller, error) {
	repo, err := repository.New(logger, cfg.DBType, cfg.DBDNS)
	if err != nil {
		return nil, err
	}

	httpClientMetrics, err := cMetrics.NewHTTPClientMetrics(cMetrics.DefaultRegistry, cMetrics.DefaultTimeBuckets)
	if err != nil {
		return nil, err
	}

	authClient, err := authclient.New(logger, httpClientMetrics, cfg.AuthAddr, cfg.AuthToken, cMetrics.InstanceName)
	if err != nil {
		return nil, err
	}

	buttonService := button.New(repo)
	userSevice := user.New(authClient)

	return &Controller{
		addr:  cfg.APIAddr,
		debug: cfg.Debug,

		logger: logger,

		buttonService: buttonService,
		userSevice:    userSevice,
	}, nil
}

func (cnt Controller) Serve(ctx context.Context) error {
	e := echo.New()
	e.Validator = vldr{validator: validator.New()}

	// FIXME: нужные мидлвари
	// e.Use()

	e.POST("/api/v1/user", cnt.createUser)
	e.DELETE("/api/v1/user", cnt.logout)
	e.GET("/api/v1/user", cnt.getUser)

	e.POST("/api/v1/button", cnt.pressButton)
	e.GET("/api/v1/button", cnt.buttons)
	e.GET("/api/v1/button/power", cnt.buttonPower)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.JSON(http.StatusInternalServerError, errorModel{
			Message: err.Error(),
		})
	}

	go func() {
		<-ctx.Done()
		err := e.Shutdown(context.TODO())
		if err != nil {
			cnt.logger.Error("shutdown http", slog.Any("error", err))
		}
	}()

	err := e.Start(cnt.addr)
	if err != nil {
		return err
	}

	return nil
}
