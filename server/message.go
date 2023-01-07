package server

import (
	"reptile-go/model"
	"time"
)

type MessageService struct {
}

// 添加消息
func (service *MessageService) AddMessage(
	msg model.ChatDetail) error {
	_, err := DbEngin.InsertOne(model.ChatDetail{
		ChatType:   msg.ChatType,
		CreateTime: time.Now().Unix() * 1000,
		FromId:     msg.FromId,
		FromName:   msg.FromName,
		FromAvatar: msg.FromAvatar,
		Data:       msg.Data,
		IsRemove:   msg.IsRemove,
		MsgType:    msg.MsgType,
		IsSend:     msg.IsSend,
		SendStatus: msg.SendStatus,
		ToAvatar:   msg.ToAvatar,
		ToId:       msg.ToId,
		ToName:     msg.ToName,
		ToUid:      msg.ToUid,
		Type:       msg.Type,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加群聊消息
func (service *MessageService) AddGroupMessage(
	msg model.ChatDetail) error {
	_, err := DbEngin.InsertOne(model.ChatGroupDetail{
		ChatType:   msg.ChatType,
		CreateTime: time.Now().Unix() * 1000,
		FromId:     msg.FromId,
		FromName:   msg.FromName,
		FromAvatar: msg.FromAvatar,
		Data:       msg.Data,
		IsRemove:   msg.IsRemove,
		MsgType:    msg.MsgType,
		IsSend:     msg.IsSend,
		SendStatus: msg.SendStatus,
		ToAvatar:   msg.ToAvatar,
		ToId:       msg.ToId,
		ToName:     msg.ToName,
		ToUid:      msg.ToUid,
		Type:       msg.Type,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加列表数据
func (service *MessageService) AddMessageList(msg []model.Message) error {
	_, err := DbEngin.Insert(&msg)
	if err != nil {
		return err
	}
	return nil
}

//获取聊天记录
func (service *MessageService) GetChatHistory(userId, dstId int64, cmd string, pageForm, pageSize int) []model.Message {
	message := make([]model.Message, 0)
	if cmd == "login" {
		DbEngin.Where("dstid = ? and cmd = ?", dstId, cmd).Desc("id").Limit(pageSize, pageForm).Find(&message)
		return message
	}
	DbEngin.Where("(userid = ? and dstid = ?) or (dstid = ? and userid = ?)", userId, dstId, userId, dstId).
		Desc("id").And("cmd = ?", cmd).Limit(pageSize, pageForm).Find(&message)
	return message
}

// 获取用户未阅读的历史消息
func (service *MessageService) GetUserNoReadHistory(userId int64) []model.ChatDetail {
	//	TODO 获取用户全部群ID
	msg := make([]model.ChatDetail, 0)
	DbEngin.Where("to_uid = ? and send_status = ?", userId, "pending").Find(&msg)

	return msg
}

// 更新用户历史消息已读回执
func (service *MessageService) UpdateNoReadHistory(ids []string) {
	tmp := model.ChatDetail{}
	tmp.SendStatus = "success"
	DbEngin.In("id", ids).Cols("send_status").Update(&tmp)

}

// 获取用户未阅读的群组历史消息
func (service *MessageService) GetNoReadGroupHistory(userId int64) []model.GroupUser {
	//	TODO 获取用户全部群ID
	msg := make([]model.GroupUser, 0)
	DbEngin.Where("user_id = ? ", userId).Find(&msg)

	return msg
}

// 更新用户历史消息已读回执
func (service *MessageService) UpdateNoReadGroupHistory(ids []string) {
	tmp := model.ChatGroupDetail{}
	tmp.SendStatus = "success"
	DbEngin.In("id", ids).Cols("send_status").Update(&tmp)

}
