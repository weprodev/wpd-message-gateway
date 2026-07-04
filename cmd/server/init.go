package main

import (
	"fmt"

	pkglogger "github.com/weprodev/go-pkg/logger"

	"github.com/weprodev/wpd-message-gateway/internal/app"
	applogger "github.com/weprodev/wpd-message-gateway/internal/infrastructure/logger"
)

func loadAppConfig(configPath string) (*app.Config, error) {
	return app.LoadAppConfig(configPath)
}

func validateAppConfig(cfg *app.Config) error {
	if err := app.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("configuration error: %w (hint: copy configs/local.example.yml to configs/local.yml)", err)
	}
	return nil
}

func setLogger(cfg *app.Config) (*pkglogger.Logger, error) {
	return applogger.New(cfg.EnvironmentName())
}

func wireApplication(cfg *app.Config) (*app.Application, error) {
	sysLogger, err := setLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	return app.Wire(cfg, sysLogger)
}
