package roll

import (
	"context"
	"dndroller/internal/logger"
	"dndroller/internal/model"
	"dndroller/internal/repo"
	"dndroller/internal/service/common"
	"fmt"
	"math"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Roll struct {
	logger *logger.Zap
	repo   *repo.Repo

	rm *rollMap
}

func NewRoll(
	log *logger.Zap,
	repo *repo.Repo,
) *Roll {
	return &Roll{
		logger: log,
		repo:   repo,
		rm:     newRollMap(),
	}
}
func (x *Roll) InlineRolls(msd *model.MsgData) (tgbotapi.Chattable, error) {
	if msd.MsgId == nil && msd.Action == nil {
		return x.makeDefaultReply(msd)
	} else {
		return x.makeRollReply(msd)
	}
}

func (x *Roll) makeDefaultReply(msd *model.MsgData) (tgbotapi.Chattable, error) {
	keyboard, err := x.setsButtons(msd.UserId)
	if err != nil {
		return nil, x.logger.LogAndWrapError(err, "")
	}
	if msd.MsgId == nil {
		msg := tgbotapi.NewMessage(msd.ChatId, "Sets")
		if keyboard != nil {
			msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
		}
		return msg, nil
	} else {
		msg := tgbotapi.NewEditMessageText(msd.ChatId, *msd.MsgId, "Sets")
		if keyboard != nil {
			msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
		}
		return msg, nil
	}
}

func (x *Roll) makeRollReply(msd *model.MsgData) (tgbotapi.Chattable, error) {
	roll, err := x.getRoll(msd)
	if err != nil {
		return nil, err
	}
	emt := tgbotapi.NewEditMessageText(msd.ChatId, *msd.MsgId, roll.ToMarkdown())
	emt.ParseMode = "MarkdownV2"
	keyboard, err := x.setsButtons(msd.UserId)
	if err != nil {
		return nil, x.logger.LogAndWrapError(err, "")
	}
	rerollKeyboard := x.rerollButtons(msd.UserId, roll)
	if keyboard != nil {
		rerollKeyboard = append(rerollKeyboard, keyboard...)
	}
	emt.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rerollKeyboard}
	return emt, err
}

func (x *Roll) getRoll(msd *model.MsgData) (*model.Roll, error) {
	var roll *model.Roll = nil
	switch msd.Action.Type {
	case model.Action_SetRoll:
		set, err := x.repo.Set.Get(context.Background(), msd.Action.Id)
		if err != nil {
			return nil, err
		}
		roll = set.Data.Roll()
		x.rm.StoreWithTll(msd.UserId, roll, time.Second*300)
	case model.Action_SetRerollDice:
		var ok bool
		roll, ok = x.rm.Load(msd.UserId)
		if !ok {
			return nil, fmt.Errorf("failed to load roll")
		}
		roll.RerollDice(msd.Action.Ind[0], msd.Action.Ind[1])
	case model.Action_SetReroll:
		var ok bool
		roll, ok = x.rm.Load(msd.UserId)
		if !ok {
			return nil, fmt.Errorf("failed to load roll")
		}
		roll.Reroll()
	}
	return roll, nil
}

func (x *Roll) setsButtons(id int) ([][]tgbotapi.InlineKeyboardButton, error) {
	sets, err := common.GetSets(x.repo, int(id))
	if err != nil {
		return nil, x.logger.LogAndWrapError(err, "")
	}
	if len(sets) == 0 {
		return nil, nil
	}
	res := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, s := range sets {
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(s.Data.String(), model.NewAction(model.Action_SetRoll, s.ID).Str()),
		}
		res = append(res, row)
	}

	return res, nil
}

func (x *Roll) rerollButtons(id int, roll *model.Roll) [][]tgbotapi.InlineKeyboardButton {
	res := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i, rr := range roll.Result {
		if len(rr) > 8 {
			res = append(res, x.rerollButtonsAligned(id, i, roll)...)
		} else {
			row := make([]tgbotapi.InlineKeyboardButton, 0, len(rr))
			for j := range rr {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(roll.FormatDice(i, j), model.NewActionInd(model.Action_SetRerollDice, id, i, j).Str()))
			}
			res = append(res, row)
		}
	}
	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Reroll", model.NewAction(model.Action_SetReroll, id).Str()),
	}
	res = append(res, row)
	return res
}

func (x *Roll) rerollButtonsAligned(id, i int, roll *model.Roll) [][]tgbotapi.InlineKeyboardButton {
	l := len(roll.Result[i])
	numRows := int(math.Ceil(float64(l) / 8))
	baseLen := l / numRows
	remainder := l % numRows
	lo := make([][]tgbotapi.InlineKeyboardButton, 0, numRows)
	j := 0
	for r := 0; r < numRows; r++ {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		rowLen := baseLen
		if remainder > 0 {
			rowLen++
			remainder--
		}
		cap := j + rowLen

		for ; j < cap; j++ {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(roll.FormatDice(i, j), model.NewActionInd(model.Action_SetRerollDice, id, i, j).Str()))
		}

		lo = append(lo, row)
	}
	return lo
}
