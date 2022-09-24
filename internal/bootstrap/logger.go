// Package bootstrap
package bootstrap

import (
	"gitlab.com/go_graphql/internal/appctx"
	"gitlab.com/go_graphql/pkg/logger"
	"gitlab.com/go_graphql/pkg/util"
)

func RegistryLogger(cfg *appctx.Config) {
	logger.Setup(logger.Config{
		Environment: util.EnvironmentTransform(cfg.App.Env),
		Debug:       cfg.App.Debug,
		Level:       cfg.Logger.Level,
		ServiceName: cfg.Logger.Name,
	})
}
