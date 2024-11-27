package o11y

import (
	"context"
	"os"

	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type NewRelicConfig struct {
	ServiceName string `mapstructure:"serviceName"`
	LicenseKey  string `mapstructure:"licenseKey"`
}

func NewObservabilityProvider(ctx context.Context, cfg *NewRelicConfig, log logger.ILogger) (*newrelic.Application, error) {
	log.Infof("Starting NR license ky : {%v}", cfg.LicenseKey)
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.ServiceName),
		newrelic.ConfigLicense(cfg.LicenseKey),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		log.Errorf("Error starting New Relic instance : {%v}", err)
	}

	return app, err

}
