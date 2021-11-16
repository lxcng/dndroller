package service

import (
	"dndroller/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (x *Bot) handleUpdate(update tgbotapi.Update) {
	x.logUpdate(update)
	var err error
	var msg tgbotapi.Chattable
	if update.Message != nil {
		msg, err = x.handleMessage(update.Message)
	}
	if update.CallbackQuery != nil {
		msg, err = x.handleCallback(update.CallbackQuery)
	}
	if err != nil {
		x.logger.Error(err)
	} else if msg != nil {
		_, err = x.bot.Send(msg)
		if err != nil {
			x.logger.Error(err)
		}
	}
}

func (x *Bot) handleMessage(msg *tgbotapi.Message) (tgbotapi.Chattable, error) {
	command := msg.Command()
	msd := model.NewMsgData().
		SetChatId(msg.Chat.ID).
		SetUserId(msg.From.ID)
	switch command {
	case "sets":
		return x.settings.InlineSettings(msd)
	case "roll":
		return x.roll.InlineRolls(msd)
	}
	return nil, nil
}

func (x *Bot) handleCallback(cb *tgbotapi.CallbackQuery) (resp tgbotapi.Chattable, err error) {
	defer x.answerCallback(cb)
	comm := model.NewActionStr(cb.Data)
	msd := model.NewMsgData().
		SetChatId(cb.Message.Chat.ID).
		SetUserId(cb.From.ID).
		SetAction(comm).
		SetMsgId(cb.Message.MessageID)
	switch comm.Type {
	case model.Action_SetRoll, model.Action_SetReroll, model.Action_SetRerollDice:
		resp, err = x.roll.InlineRolls(msd)
	default:
		resp, err = x.settings.InlineSettings(msd)
	}
	return
}

func (x *Bot) answerCallback(cb *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(cb.ID, "")
	callback.ShowAlert = false
	_, err := x.bot.AnswerCallbackQuery(callback)
	if err != nil {
		x.logger.Errorw("failed to answer callback query", "err", err)
	}
}

func (x *Bot) logUpdate(update tgbotapi.Update) {
	agrs := []interface{}{}
	switch {
	case update.Message != nil:
		agrs = append(agrs, "message", update.Message)
	case update.CallbackQuery != nil:
		agrs = append(agrs, "callbackQuery", update.CallbackQuery)
	}
	x.logger.Infow("update", agrs...)
}
