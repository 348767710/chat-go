package server

import (
	"errors"
	"fmt"
	"math/rand"
	"reptile-go/model"
	"reptile-go/util"
	"time"

	"github.com/prometheus/common/log"
)

type UserService struct {
}

// 注册
func (s *UserService) Register(username, plainpwd, nickname, avatar, sex string) (user model.User, err error) {
	// 手机号
	// 明文密码
	// 昵称
	// 1.检测手机号是否存在
	tmp := model.User{}
	_, err = DbEngin.Where("username=? ", username).Get(&tmp)
	if err != nil {
		return tmp, err
	}
	// 如果存在则返回提示已注册
	if tmp.Id > 0 {
		return tmp, errors.New("该账号已经注册")
	}
	//否则插入数据
	tmp.Username = username
	tmp.Avatar = avatar
	tmp.Nickname = nickname
	tmp.Sex = sex
	tmp.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	tmp.Password = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.RegisterDate = time.Now()
	//头像
	tmp.Avatar = "/static/images/default_avatar.jpg"
	// token 随机数
	tmp.Token = fmt.Sprintf("%08d", rand.Int31())
	// 插入数据
	_, err = DbEngin.InsertOne(&tmp)
	log.Warn()
	//返回新用户信息
	return tmp, err
}

// 登录
func (s *UserService) Login(username, plainpwd string) (user model.User, err error) {
	tmp := model.User{}
	_, err = DbEngin.Where("username=?", username).Get(&tmp)
	if err != nil {
		return tmp, err
	}
	// 账号不存在
	if tmp.Id == 0 {
		return tmp, errors.New("账号不存在")
	}
	// 检测密码
	if !util.ValidatePasswd(plainpwd, tmp.Salt, tmp.Password) {
		return tmp, errors.New("密码不正确")
	}
	// 刷新Token
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.MD5Encode(str)
	tmp.Token = token
	tmp.Online = 1 //在线
	DbEngin.ID(tmp.Id).Cols("token,online").Update(&tmp)
	// 返回数据
	return tmp, nil
}

// 查询某个用户信息
func (s *UserService) Find(userId int64) (user model.User) {
	tmp := model.User{}
	DbEngin.ID(userId).Get(&tmp)
	// 返回数据
	return tmp
}

// 更新用户数据(头像)
func (s *UserService) UserInfo(userId int64, avatar string) {
	tmp := model.User{}
	tmp.Avatar = avatar
	DbEngin.ID(userId).Cols("avatar").Update(&tmp)
}

// 更新用户昵称
func (s *UserService) UserNickname(userId int64, nickname string) {
	tmp := model.User{}
	tmp.Nickname = nickname
	DbEngin.ID(userId).Cols("nickname").Update(&tmp)
}

// 获取用户数据
func (s *UserService) GetUser(userId int64) (user model.User, err error) {
	tmp := model.User{}
	DbEngin.Where("id=? ", userId).Get(&tmp)
	return tmp, nil
}

// 通过手机账号获取用户数据
func (s *UserService) GetUserByName(mobile int64) (user model.User, err error) {
	tmp := model.User{}
	DbEngin.Where("username=? ", mobile).Get(&tmp)
	// 账号不存在
	if tmp.Id == 0 {
		return tmp, errors.New("账号不存在")
	}
	return tmp, nil
}

// 批量获取用户数据
func (s *UserService) GetUserByIds(ids []string) []model.User {
	tmp := make([]model.User, 0)
	DbEngin.In("id", ids).Cols("id,username,nickname,avatar").Find(&tmp)
	return tmp
}

// 更新用户密码
func (s *UserService) EditUserPwd(userId int64, passwd string) {
	userinfo := s.Find(userId)
	tmp := model.User{}
	tmp.Password = util.MakePasswd(passwd, userinfo.Salt)
	DbEngin.ID(userId).Cols("password").Update(&tmp)
}

// 退出登陆
func (s *UserService) UserOnline(userId int64, online int) {
	tmp := model.User{}
	tmp.Online = online
	DbEngin.ID(userId).Cols("online").Update(&tmp)
}
