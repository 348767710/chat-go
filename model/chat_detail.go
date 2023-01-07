package model

type ChatDetail struct {
	Id         int64  `json:"id,omitempty" form:"id"`                   //消息ID
	ChatType   string `json:"chat_type,omitempty" form:"chat_type"`     //聊天类型(user,group)
	CreateTime int64  `json:"create_time,omitempty" form:"create_time"` //创建时间
	Data       string `json:"data,omitempty" form:"data"`               //消息内容
	FromAvatar string `json:"from_avatar,omitempty" form:"from_avatar"` //发送者头像
	FromId     int64  `json:"from_id,omitempty" form:"from_id"`         //发送者id
	FromName   string `json:"from_name,omitempty" form:"from_name"`     //发送者名字
	IsRemove   int64  `json:"is_remove,omitempty" form:"is_remove"`     //是否删除
	MsgType    string `json:"msg_type,omitempty" form:"msg_type"`       //消息类型(text,img)
	IsSend     string `json:"is_send,omitempty" form:"is_send"`         //
	SendStatus string `json:"send_status,omitempty" form:"send_status"` //发送状态
	ToAvatar   string `json:"to_avatar,omitempty" form:"to_avatar"`     //
	ToId       int64  `json:"to_id,omitempty" form:"to_id"`             //
	ToName     string `json:"to_name,omitempty" form:"to_name"`         //
	ToUid      int64  `json:"to_uid,omitempty" form:"to_uid"`           //
	Type       string `json:"type,omitempty" form:"type"`               // 类型
}

func (ChatDetail) TableName() string {
	//方法一：指定数据库表名称为chat_users
	return "chat_detail"
}
