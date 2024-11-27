package echoserver

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"time"
	nrecho "github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

type EchoConfig struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"`
	Host                string   `mapstructure:"host"`
}

func NewEchoServer(o11y *newrelic.Application) *echo.Echo {

	e := echo.New()
	e.Use(nrecho.Middleware(o11y))
	return e
}

func RunHttpServer(ctx context.Context, echo *echo.Echo, log logger.ILogger, cfg *EchoConfig) error {
	echo.Server.ReadTimeout = ReadTimeout
	echo.Server.WriteTimeout = WriteTimeout
	echo.Server.MaxHeaderBytes = MaxHeaderBytes

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Infof("shutting down Http PORT: {%s}", cfg.Port)
				err := echo.Shutdown(ctx)
				if err != nil {
					log.Errorf("(Shutdown) err: {%v}", err)
					return
				}
				log.Info("server exited properly")
				return
			}
		}
	}()

	err := echo.Start(cfg.Port)

	return err
}

func ApplyVersioningFromHeader(echo *echo.Echo) {
	echo.Pre(apiVersion)
}

func RegisterGroupFunc(groupName string, echo *echo.Echo, builder func(g *echo.Group)) *echo.Echo {
	builder(echo.Group(groupName))

	return echo
}

// APIVersion Header Based Versioning
func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		apiVersion := headers.Get("version")

		req.URL.Path = fmt.Sprintf("/%s%s", apiVersion, req.URL.Path)

		return next(c)
	}
}
