package main

import (
	"context"
	"dndroller/internal/config"
	"dndroller/internal/logger"
	"dndroller/internal/repo"
	"log"
	"os"
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
	cl, err := repo.NewCl(cfg, log)
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()
	f, err := os.Create("migrate.sql")
	if err != nil {
		log.Fatalf("create migrate file: %v", err)
	}
	defer f.Close()
	if err := cl.Schema.WriteTo(context.Background(), f); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
