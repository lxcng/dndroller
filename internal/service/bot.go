package service

import (
	"context"
	"dndroller/internal/config"
	"dndroller/internal/logger"
	"dndroller/internal/repo"
	"dndroller/internal/service/roll"
	"dndroller/internal/service/settings"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	logger *logger.Zap
	bot    *tgbotapi.BotAPI
	repo   *repo.Repo
	token  string

	roll     *roll.Roll
	settings *settings.Settings
}

func NewBot(cfg *config.Config, log *logger.Zap, repo *repo.Repo) *Bot {
	return &Bot{
		logger:   log,
		token:    cfg.Token,
		repo:     repo,
		roll:     roll.NewRoll(log, repo),
		settings: settings.NewSettings(log, repo),
	}
}

func (that *Bot) init() error {
	if err := tgbotapi.SetLogger(that.logger); err != nil {
		return err
	}
	var bot *tgbotapi.BotAPI
	var err error
	bot, err = tgbotapi.NewBotAPI(that.token)
	if err != nil {
		return that.logger.LogAndWrapError(err, "failed")
	}
	that.bot = bot
	return nil
}

func (that *Bot) startLoop() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := that.bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("error getting updates channel %q", err.Error())
	}

	go func() {
		for update := range updates {
			go that.handleUpdate(update)
		}
	}()

	return nil
}

func (that *Bot) Start(_ context.Context) error {
	err := that.init()
	if err != nil {
		return err
	}
	return that.startLoop()
}

func (that *Bot) Stop(_ context.Context) error {
	that.bot.StopReceivingUpdates()
	return nil
}
