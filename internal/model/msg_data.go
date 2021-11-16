package model

type MsgData struct {
	ChatId int64
	UserId int
	Action *Action
	MsgId  *int
}

func NewMsgData() *MsgData {
	return &MsgData{}
}

func (x *MsgData) SetChatId(id int64) *MsgData {
	x.ChatId = id
	return x
}

func (x *MsgData) SetUserId(id int) *MsgData {
	x.UserId = id
	return x
}

func (x *MsgData) SetMsgId(id int) *MsgData {
	x.MsgId = &id
	return x
}

func (x *MsgData) SetAction(ac *Action) *MsgData {
	x.Action = ac
	return x
}
