package factory

import (
	"fmt"
	"log"

	"github.com/manish-npx/go-student-api/internal/config"
	"github.com/manish-npx/go-student-api/internal/storage"
	"github.com/manish-npx/go-student-api/internal/storage/postgres"
	"github.com/manish-npx/go-student-api/internal/storage/sqlite"
)

// 🏭 Register each database's constructor
var factories = map[string]func(config.Config) (storage.Storage, error){
	"sqlite": func(cfg config.Config) (storage.Storage, error) {
		return sqlite.New(cfg)
	},
	"postgres": func(cfg config.Config) (storage.Storage, error) {
		return postgres.New(cfg)
	},
}

// 🧩 Main entrypoint for selecting DB
func NewStorage(cfg config.Config) (storage.Storage, error) {
	if cfg.DBType == "" {
		log.Fatal("❌ No db_driver specified in config.yaml")
	}
	createFn, ok := factories[cfg.DBType]
	if !ok {
		return nil, fmt.Errorf(
			"unsupported db type: %s (supported: sqlite, postgres)",
			cfg.DBType,
		)
	}
	return createFn(cfg)
}
