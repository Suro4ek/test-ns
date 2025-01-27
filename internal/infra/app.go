package infra

import (
	"errors"
	"test-ns/config"
	"test-ns/internal/adapters/repo"
	"test-ns/internal/infra/database"
	"test-ns/internal/infra/router"
)

type app struct {
	cfg         config.Config
	router      router.Router
	database    repo.GSQL
	transaction repo.GTransaction
}

func NewConfig(config config.Config) *app {
	return &app{
		cfg: config,
	}
}

func (a *app) DBGSql() *app {
	db := database.NewGormDB(a.cfg)
	a.database = *db
	return a
}

func (a *app) GTransaction() *app {
	if a.database == nil {
		panic(errors.New("database is not initialized"))
	}
	transaction := database.NewGormTransaction(a.database)
	a.transaction = transaction
	return a
}

func (a *app) Run() {
	a.router = router.NewRouter(a.transaction, a.database)
	a.router.Start()
}
