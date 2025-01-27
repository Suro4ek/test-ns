package main

import (
	"test-ns/config"
	"test-ns/internal/infra"
)

func main() {
	cfg, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}
	infra.NewConfig(cfg).DBGSql().GTransaction().Run()
}
