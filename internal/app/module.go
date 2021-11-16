package app

import (
	"dndroller/internal/config"
	"dndroller/internal/logger"
	"dndroller/internal/repo"
	"dndroller/internal/service"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module(cfg *config.Config, log *logger.Zap) fx.Option {
	return fx.Options(
		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(10*time.Second),
		fx.Logger(log),
		fx.Provide(
			func() *config.Config {
				return cfg
			}),

		fx.Provide(
			func() *logger.Zap {
				return log
			},
		),
		repo.Module(),
		service.Module(),

		fx.Invoke(
			func(cfg *config.Config, log *logger.Zap) {
				log.Info("Telegram bot started")
			},
		),
	)
}

func Recover(log *logger.Zap) {
	if err := recover(); err != nil {
		log.Fatal("app recover error", zap.Any("recoveryError", err))
	}
}
