package model

// 群组
type ChatGroup struct {
	Id         int64  `form:"id" json:"id"`
	Admin      string `xorm:"varchar(100)" form:"admin",json:"admin"`      // 名称
	AdminId    int64  `xorm:"varchar(20)" form:"admin_id" json:"admin_id"` // 群主ID
	Title      string `xorm:"varchar(120)" form:"title" json:"title"`      // 描述
	CreateTime int64  `xorm:"create_time" form:"create_time" json:"create_time"`
	GroupUsers []User `xorm:"-" json:"group_users"`
}

func (ChatGroup) TableName() string {
	//方法一：指定数据库表名称为chat_users
	return "chat_group"
}
