package model

// 可根据具体业务做拆分
type ChatFriends struct {
	Id       int64  `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	UserId   int64  `xorm:"int(10)" form:"user_id" json:"user_id" validate:"required|int"`         // 谁的10000
	FriendId int64  `xorm:"varchar(10)" form:"friend_id" json:"friend_id" validate:"required|int"` // 对端，10001
	PetName  string `json:"pet_name,omitempty" form:"pet_name"`                                    // 备注昵称

}

func (ChatFriends) TableName() string {
	//方法一：指定数据库表名称为chat_users
	return "chat_friends"
}
