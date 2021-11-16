package settings

import (
	"dndroller/internal/logger"
	"dndroller/internal/model"
	"dndroller/internal/repo"
	"dndroller/internal/service/common"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Settings struct {
	logger *logger.Zap
	repo   *repo.Repo
}

func NewSettings(
	log *logger.Zap,
	repo *repo.Repo,
) *Settings {
	return &Settings{
		logger: log,
		repo:   repo,
	}
}

func (x *Settings) InlineSettings(msd *model.MsgData) (tgbotapi.Chattable, error) {
	if msd.Action == nil {
		return x.makeDefaultReply(msd, "")
	}
	switch msd.Action.Type {
	case model.Action_Settings, model.Action_Back:
		return x.makeDefaultReply(msd, "")
	case model.Action_SetDelete:
		set, err := common.DeleteSet(x.repo, msd.Action.Id)
		if err != nil {
			return nil, err
		}
		return x.makeDefaultReply(msd, fmt.Sprintf("Deleted set: %s", set.Data.String()))
	default:
		return x.makeSetReply(msd, "")
	}
}

func (x *Settings) makeDefaultReply(msd *model.MsgData, title string) (tgbotapi.Chattable, error) {
	if title == "" {
		title = "Settings"
	}
	keyboard, err := x.setsButtons(msd.UserId)
	if err != nil {
		return nil, x.logger.LogAndWrapError(err, "")
	}
	if msd.MsgId == nil {
		msg := tgbotapi.NewMessage(msd.ChatId, title)
		if keyboard != nil {
			msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
		}
		return msg, nil
	} else {
		msg := tgbotapi.NewEditMessageText(msd.ChatId, *msd.MsgId, title)
		if keyboard != nil {
			msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
		}
		return msg, err
	}
}

func (x *Settings) setsButtons(id int) ([][]tgbotapi.InlineKeyboardButton, error) {
	sets, err := common.GetSets(x.repo, int(id))
	if err != nil {
		return nil, x.logger.LogAndWrapError(err, "")
	}

	res := make([][]tgbotapi.InlineKeyboardButton, 0, len(sets)+1)
	for _, s := range sets {
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(s.Data.String(), model.NewAction(model.Action_SetEdit, s.ID).Str()),
		}
		res = append(res, row)
	}
	res = append(res, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Add", model.NewAction(model.Action_AddSet, id).Str()),
	})

	return res, nil
}
