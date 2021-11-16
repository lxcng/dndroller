package main

import (
	"dndroller/internal/app"
	"dndroller/internal/config"
	"dndroller/internal/logger"
	"log"

	"go.uber.org/fx"
)

func main() {
	cfg, err := config.NewConfigEnv()
	if err != nil {
		log.Fatal(err)
	}
	log, err := logger.NewZapLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Recover(log)
	fx.New(
		app.Module(cfg, log),
	).Run()

}
