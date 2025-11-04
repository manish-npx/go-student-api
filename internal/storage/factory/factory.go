package factory

import (
	"fmt"

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
	createFn, ok := factories[cfg.DBType]
	if !ok {
		return nil, fmt.Errorf("unsupported db type: %s (use sqlite or postgres)", cfg.DBType)
	}
	return createFn(cfg)
}
