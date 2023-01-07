package model

// 群成员表
type ChatGroupDetailIds struct {
	Id         int64  `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	CreateTime int64  `xorm:"create_time" form:"create_time" json:"create_time"` //
	SendStatus string `json:"send_status,omitempty" form:"send_status"`          //发送状态
	UserId     int64  `xorm:"int(10)" form:"user_id" json:"user_id"`             // 会员id
}

func (ChatGroupDetailIds) TableName() string {
	//方法一：指定数据库表名称为chat_users
	return "chat_group_detail_ids"
}
