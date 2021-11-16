package settings

import (
	"context"
	"dndroller/internal/model"
	"dndroller/internal/repo"
	"dndroller/internal/service/common"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (x *Settings) makeSetReply(msd *model.MsgData, title string) (tgbotapi.Chattable, error) {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	var set *repo.Set
	var err error
	switch msd.Action.Type {
	case model.Action_AddSet:
		set, err = common.CreateSet(x.repo, msd.Action.Id)
		if err != nil {
			return nil, err
		}
	case model.Action_SetEdit:
		set, err = x.repo.Set.Get(context.Background(), msd.Action.Id)
		if err != nil {
			return nil, err
		}
	default:
		set, err = x.mutateSet(msd)
		if err != nil {
			return nil, err
		}
	}
	buttons = append(buttons, x.diceButtons(set.ID)...)
	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Delete", model.NewAction(model.Action_SetDelete, set.ID).Str()),
		tgbotapi.NewInlineKeyboardButtonData("Back", model.NewAction(model.Action_Back, msd.UserId).Str()),
	}
	buttons = append(buttons, row)
	emt := tgbotapi.NewEditMessageText(msd.ChatId, *msd.MsgId, fmt.Sprintf("Set %s", set.Data.String()))
	emt.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}
	return emt, nil
}

func (x *Settings) mutateSet(msd *model.MsgData) (*repo.Set, error) {
	ctx := context.Background()
	set, err := x.repo.Set.Get(ctx, msd.Action.Id)
	if err != nil {
		return nil, err
	}
	switch msd.Action.Type {
	case model.Action_DiceD4M:
		set.Data.Sub(4)
	case model.Action_DiceD4P:
		set.Data.Add(4)

	case model.Action_DiceD6M:
		set.Data.Sub(6)
	case model.Action_DiceD6P:
		set.Data.Add(6)

	case model.Action_DiceD8M:
		set.Data.Sub(8)
	case model.Action_DiceD8P:
		set.Data.Add(8)

	case model.Action_DiceD10M:
		set.Data.Sub(10)
	case model.Action_DiceD10P:
		set.Data.Add(10)

	case model.Action_DiceD12M:
		set.Data.Sub(12)
	case model.Action_DiceD12P:
		set.Data.Add(12)

	case model.Action_DiceD20M:
		set.Data.Sub(20)
	case model.Action_DiceD20P:
		set.Data.Add(20)

	case model.Action_DiceD100M:
		set.Data.Sub(100)
	case model.Action_DiceD100P:
		set.Data.Add(100)
	}

	return x.repo.Set.UpdateOne(set).SetData(set.Data).Save(ctx)
	// return x.repo.Set.UpdateOne(set).Save(ctx)
}

func (x *Settings) diceButtons(id int) [][]tgbotapi.InlineKeyboardButton {
	res := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("-d4", model.NewAction(model.Action_DiceD4M, id).Str()),
			tgbotapi.NewInlineKeyboardButtonData("+d4", model.NewAction(model.Action_DiceD4P, id).Str()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("-d6", model.NewAction(model.Action_DiceD6M, id).Str()),
			tgbotapi.NewInlineKeyboardButtonData("+d6", model.NewAction(model.Action_DiceD6P, id).Str()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("-d8", model.NewAction(model.Action_DiceD8M, id).Str()),
			tgbotapi.NewInlineKeyboardButtonData("+d8", model.NewAction(model.Action_DiceD8P, id).Str()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("-d10", model.NewAction(model.Action_DiceD10M, id).Str()),
			tgbotapi.NewInlineKeyboardButtonData("+d10", model.NewAction(model.Action_DiceD10P, id).Str()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("-d12", model.NewAction(model.Action_DiceD12M, id).Str()),
			tgbotapi.NewInlineKeyboardButtonData("+d12", model.NewAction(model.Action_DiceD12P, id).Str()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("-d20", model.NewAction(model.Action_DiceD20M, id).Str()),
			tgbotapi.NewInlineKeyboardButtonData("+d20", model.NewAction(model.Action_DiceD20P, id).Str()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("-d100", model.NewAction(model.Action_DiceD100M, id).Str()),
			tgbotapi.NewInlineKeyboardButtonData("+d100", model.NewAction(model.Action_DiceD100P, id).Str()),
		},
	}
	return res
}
