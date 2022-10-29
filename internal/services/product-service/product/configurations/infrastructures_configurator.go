package configurations

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http_client"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	consumers2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/consumers"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	events3 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/events"
	events2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type infrastructureConfigurator struct {
	Log  logger.ILogger
	Cfg  *config.Config
	Echo *echo.Echo
}

func NewInfrastructureConfigurator(log logger.ILogger, cfg *config.Config, echo *echo.Echo) *infrastructureConfigurator {
	return &infrastructureConfigurator{Cfg: cfg, Echo: echo, Log: log}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (error, func()) {

	infrastructure := &shared.InfrastructureConfiguration{Cfg: ic.Cfg, Echo: ic.Echo, Log: ic.Log, Validator: validator.New()}

	cleanups := []func(){}

	gorm, err := ic.configGorm()
	if err != nil {
		return err, nil
	}
	infrastructure.Gorm = gorm

	tp, err := ic.configOpenTelemetry(ctx)
	if err != nil {
		return err, nil
	}
	infrastructure.JaegerTracer = tp.Tracer(ic.Cfg.Jaeger.TracerName)

	ic.Log.Infof("%s is running", config.GetMicroserviceName(ic.Cfg.ServiceName))

	conn, err, rabbitMqCleanup := rabbitmq.NewRabbitMQConn(ic.Cfg.Rabbitmq)
	if err != nil {
		return err, nil
	}

	infrastructure.ConnRabbitmq = conn
	cleanups = append(cleanups, rabbitMqCleanup)

	infrastructure.RabbitmqPublisher = rabbitmq.NewPublisher(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer)

	createProductConsumer := rabbitmq.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers2.HandleConsumeCreateProduct)
	updateProductConsumer := rabbitmq.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers2.HandleConsumeUpdateProduct)
	deleteProductConsumer := rabbitmq.NewConsumer(ic.Cfg.Rabbitmq, infrastructure.ConnRabbitmq, infrastructure.Log, infrastructure.JaegerTracer, consumers2.HandleConsumeDeleteProduct)

	go func() {
		err = createProductConsumer.ConsumeMessage(ctx, v1.ProductCreated{})
		if err != nil {
			ic.Log.Error(err)
		}
	}()

	go func() {
		err = updateProductConsumer.ConsumeMessage(ctx, events2.ProductUpdated{})
		if err != nil {
			ic.Log.Error(err)
		}
	}()

	go func() {
		err = deleteProductConsumer.ConsumeMessage(ctx, events3.ProductDeleted{})
		if err != nil {
			ic.Log.Error(err)
		}
	}()

	httpClient := http_client.NewHttpClient()
	infrastructure.HttpClient = httpClient

	ic.configSwagger()

	ic.configMiddlewares(ic.Cfg.Jaeger)
	if err != nil {
		return err, nil
	}

	pc := NewProductsModuleConfigurator(infrastructure)

	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return err, nil
	}

	ic.Echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", config.GetMicroserviceName(ic.Cfg.ServiceName)))
	})

	return nil, func() {
		for _, c := range cleanups {
			c()
		}
	}
}