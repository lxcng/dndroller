package model

import "encoding/json"

type actionType string

const (
	// to settings
	Action_Settings  actionType = "set_settings"
	Action_Back      actionType = "back"
	Action_SetDelete actionType = "delete"

	// to set
	Action_AddSet  actionType = "add"
	Action_SetEdit actionType = "edit_set"

	Action_DiceD4M   actionType = "d4m"
	Action_DiceD4P   actionType = "d4p"
	Action_DiceD6M   actionType = "d6m"
	Action_DiceD6P   actionType = "d6p"
	Action_DiceD8M   actionType = "d8m"
	Action_DiceD8P   actionType = "d8p"
	Action_DiceD10M  actionType = "d10m"
	Action_DiceD10P  actionType = "d10p"
	Action_DiceD12M  actionType = "d12m"
	Action_DiceD12P  actionType = "d12p"
	Action_DiceD20M  actionType = "d20m"
	Action_DiceD20P  actionType = "d20p"
	Action_DiceD100M actionType = "d100m"
	Action_DiceD100P actionType = "d100p"

	//
	Action_SetRoll       actionType = "roll"
	Action_SetRerollDice actionType = "reroll_dice"
	Action_SetReroll     actionType = "reroll"
)

type Action struct {
	Type actionType
	Id   int
	Ind  *[2]int
}

func NewAction(ac actionType, id int) *Action {
	return &Action{
		Type: ac,
		Id:   id,
	}
}

func NewActionInd(ac actionType, id, i, j int) *Action {
	return &Action{
		Type: ac,
		Id:   id,
		Ind:  &[2]int{i, j},
	}
}

func NewActionStr(str string) *Action {
	res := &Action{}
	_ = json.Unmarshal([]byte(str), res)
	return res
}

func (x *Action) Str() string {
	bt, _ := json.Marshal(x)
	return string(bt)
}
