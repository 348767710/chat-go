package model

import "time"

const (
	SEX_WOMEN  = "W" // 女
	SEX_MAN    = "M" // 男
	SEX_UNKNOW = "U" // 未知
)

type User struct {
	Id            int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`             // 用户ID
	Username      string    `xorm:"varchar(20)" form:"username" json:"username"`            // 手机号
	Password      string    `xorm:"varchar(40)" form:"password" json:"-"`                   // 用户密码 = f(plainpwd + salt),MD5,
	Avatar        string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`               // 头像
	Sex           string    `xorm:"varchar(2)" form:"sex" json:"sex"`                       // 性别
	Nickname      string    `xorm:"varchar(20)" form:"nickname" json:"nickname"`            // 昵称
	Salt          string    `xorm:"varchar(10)" form:"salt" json:"-"`                       // 加盐随机字符串6
	Online        int       `xorm:"int(10)" form:"online" json:"online"`                    // 是否在线
	Token         string    `xorm:"varchar(40)" form:"token" json:"token"`                  // Token
	Memo          string    `xorm:"varchar(140)" form:"memo" json:"memo"`                   // 备注
	RegisterDate  time.Time `xorm:"datetime" form:"register_date" json:"register_date"`     // 创建时间
	LastLoginDate time.Time `xorm:"datetime" form:"last_login_date" json:"last_login_date"` // 最后登陆时间
}

func (User) TableName() string {
	//方法一：指定数据库表名称为chat_users
	return "chat_users"
}
