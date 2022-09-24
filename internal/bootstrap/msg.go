// Package bootstrap
package bootstrap

import (
	"gitlab.com/go_graphql/internal/consts"
	"gitlab.com/go_graphql/pkg/logger"
	"gitlab.com/go_graphql/pkg/msgx"
)

func RegistryMessage() {
	err := msgx.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
