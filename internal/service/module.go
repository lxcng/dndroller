package service

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewBot),
		fx.Invoke(func(lc fx.Lifecycle, bot *Bot) {
			lc.Append(fx.Hook{
				OnStart: bot.Start,
				OnStop:  bot.Stop,
			})
		}),
	)
}
