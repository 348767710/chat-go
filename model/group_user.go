package model

// 群成员表
type GroupUser struct {
	Id      int64 `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	GroupId int64 `xorm:"int(10)" form:"group_id" json:"group_id"` // 群id
	UserId  int64 `xorm:"int(10)" form:"user_id" json:"user_id"`   // 会员id
}

func (GroupUser) TableName() string {
	//方法一：指定数据库表名称为chat_users
	return "chat_group_users"
}