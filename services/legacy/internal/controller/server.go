package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	coreLogger "github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/services/legacy/internal/repository"
	"github.com/gbh007/buttoners/services/legacy/internal/service/button"
	"github.com/gbh007/buttoners/services/legacy/internal/service/user"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

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

	kafkaTaskConsumer *kafka.Consumer[dto.KafkaTaskData]
	kafkaLogClient    *kafka.Producer[dto.KafkaLogData]
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

	Kafka struct {
		TaskTopic string `envconfig:"KAFKA_TASK_TOPIC" default:"gate"`
		LogTopic  string `envconfig:"KAFKA_LOG_TOPIC" default:"log"`
		GroupID   string `envconfig:"KAFKA_GROUP_ID"`
		Addr      string `envconfig:"KAFKA_ADDR" default:"kafka:9092"`
	}
}

func New(logger *slog.Logger, cfg Config) (*Controller, error) {
	tracer := otel.GetTracerProvider().Tracer(cMetrics.InstanceName)

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

	queueReaderMetrics, err := cMetrics.NewQueueReaderMetrics(cMetrics.DefaultRegistry, cMetrics.DefaultTimeBuckets)
	if err != nil {
		return nil, err
	}

	queueWriterMetrics, err := cMetrics.NewQueueWriterMetrics(cMetrics.DefaultRegistry, cMetrics.DefaultTimeBuckets)
	if err != nil {
		return nil, err
	}

	kafkaTaskClient := kafka.NewProducer[dto.KafkaTaskData](logger, cfg.Kafka.Addr, cfg.Kafka.TaskTopic, queueWriterMetrics)
	kafkaLogClient := kafka.NewProducer[dto.KafkaLogData](logger, cfg.Kafka.Addr, cfg.Kafka.LogTopic, queueWriterMetrics)

	buttonService := button.New(repo, kafkaTaskClient)
	userSevice := user.New(authClient)

	kafkaTaskConsumer := kafka.NewConsumer(
		logger,
		cfg.Kafka.Addr,
		cfg.Kafka.TaskTopic,
		cfg.Kafka.GroupID,
		queueReaderMetrics,
		func(ctx context.Context, key string, data dto.KafkaTaskData) error {
			ctx, span := tracer.Start(ctx, "handle button")
			defer span.End()

			ctx, rabbitCnl := context.WithTimeout(ctx, time.Second*10)
			defer rabbitCnl()

			err := buttonService.ConsumePressButton(ctx, int(data.UserID))
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "handle error")

				coreLogger.LogWithMeta(logger, ctx, slog.LevelError, "write to rabbitmq", "error", err.Error(), "msg_key", key)

				return err
			}

			return nil
		},
	)

	return &Controller{
		addr:  cfg.APIAddr,
		debug: cfg.Debug,

		logger: logger,

		buttonService:     buttonService,
		userSevice:        userSevice,
		kafkaTaskConsumer: kafkaTaskConsumer,
		kafkaLogClient:    kafkaLogClient,
	}, nil
}

func (cnt Controller) Serve(ctx context.Context) error {
	e := echo.New()
	e.Validator = vldr{validator: validator.New()}

	// FIXME: нужные мидлвари
	e.Use(
		cnt.logActivity(),
	)

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

	go cnt.kafkaTaskConsumer.Start(ctx)

	err := e.Start(cnt.addr)
	if err != nil {
		return err
	}

	err = cnt.kafkaTaskConsumer.Close()
	if err != nil {
		return err
	}

	return nil
}
