// Package migration
package migration

import (
	"gitlab.privy.id/go_graphql/pkg/databasex"
	"time"

	"gitlab.privy.id/go_graphql/internal/appctx"
)

func MigrateDatabase() {
	cfg := appctx.NewConfig()

	databasex.DatabaseMigration(&databasex.Config{
		Driver:       cfg.WriteDB.Driver,
		Host:         cfg.WriteDB.Host,
		Port:         cfg.WriteDB.Port,
		Name:         cfg.WriteDB.Name,
		User:         cfg.WriteDB.User,
		Password:     cfg.WriteDB.Pass,
		Charset:      cfg.WriteDB.Charset,
		Timeout:      cfg.WriteDB.TimeoutSecond,
		MaxIdleConns: cfg.WriteDB.MaxIdle,
		MaxOpenConns: cfg.WriteDB.MaxOpen,
		MaxLifetime:  time.Duration(cfg.WriteDB.MaxLifeTimeMS) * time.Millisecond,
		TimeZone:     cfg.WriteDB.Timezone,
	})
}
