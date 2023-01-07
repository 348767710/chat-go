package server

import (
	"errors"
	"fmt"
	"reptile-go/args"
	"reptile-go/model"
	"strconv"
	"time"
)

type ChatFriendsService struct {
}

// 添加好友
func (service *ChatFriendsService) AddFriend(userid, friend_id int64) error {
	// 如果添加自己为好友
	if userid == friend_id {
		return errors.New("无法添加自己为好友")
	}
	// 判断是否已是好友
	tmp := model.ChatFriends{}
	tmpUser := model.User{}
	// 判断将要添加的好友是否存在
	DbEngin.Where("id = ?", friend_id).Get(&tmpUser)
	// 用户不存在
	if tmpUser.Id == 0 {
		return errors.New("用户不存在!")
	}
	// 1.查询是否已经是好友了
	// 这里是条件链式操作
	// 获取1条数据
	DbEngin.Where("user_id = ?", userid).
		And("friend_id = ?", friend_id).
		Get(&tmp)
	if tmp.Id > 0 {
		return errors.New("请勿重复添加好友!")
	}
	// 开启事务
	session := DbEngin.NewSession()
	session.Begin()
	// 插入自己的好友数据
	_, e1 := session.InsertOne(model.ChatFriends{
		UserId:   userid,
		FriendId: friend_id,
	})
	// 插入对方的数据
	_, e2 := session.InsertOne(model.ChatFriends{
		UserId:   friend_id,
		FriendId: userid,
	})
	// 如果没有错误
	if e1 == nil && e2 == nil {
		// 提交事务
		session.Commit()
		return nil
	} else {
		// 回滚事务
		session.Rollback()
		if e1 != nil {
			return e1
		} else {
			return e2
		}
	}
}


// 删除好友
func (service *ChatFriendsService) DelFriend(userid, friend_id int64) error {

	// 开启事务
	session := DbEngin.NewSession()
	session.Begin()
	// 删除自己的好友数据
	_, e1 := session.Delete(model.ChatFriends{
		UserId:   userid,
		FriendId: friend_id,
	})
	// 删除对方的数据
	_, e2 := session.Delete(model.ChatFriends{
		UserId:   friend_id,
		FriendId: userid,
	})
	// 如果没有错误
	if e1 == nil && e2 == nil {
		// 提交事务
		session.Commit()
		return nil
	} else {
		// 回滚事务
		session.Rollback()
		if e1 != nil {
			return e1
		} else {
			return e2
		}
	}
}

//查找好友列表
func (service *ChatFriendsService) SearchFriend(userId int64) []model.ChatFriends {
	conconts := make([]model.ChatFriends, 0)
	// 查询好友列表
	DbEngin.Where("user_id = ? ", userId).Find(&conconts)

	return conconts
}

// 创建群
func (service *ChatFriendsService) CreateCommunity(comm model.ChatGroup) (ret model.ChatGroup, err error) {

	com := model.ChatGroup{
		AdminId: comm.AdminId,
	}
	// 判断建群数量
	num, err := DbEngin.Count(&com)
	if num > 5 {
		return ret, errors.New("一个用户最多创建5个群")
	} else {
		comm.CreateTime = time.Now().Unix() * 1000
		//添加群表
		_, err = DbEngin.InsertOne(&comm)

		return comm, err
	}
}

//获取群列表
func (service *ChatFriendsService) SearchComunity(userId int64) []model.GroupUser {
	conconts := make([]model.GroupUser, 0)
	//comIds := make([]int64, 0)
	DbEngin.Where("user_id = ? ", userId).Find(&conconts)
	return conconts
}

// 添加群成员
func (service *ChatFriendsService) JoinCommunity(arg []args.GroupArg) error {

	for _, v := range arg {
		cot := model.GroupUser{
			UserId:  v.UserId,
			GroupId: v.GroupId,
		}
		_, err := DbEngin.InsertOne(&cot)
		if err != nil {
			return err
		}

	}
	return nil

}

// 获取群ids
func (service *ChatFriendsService) SearchComunityIds(userId int64) (comIds []int64) {
	//	TODO 获取用户全部群ID
	conconts := make([]model.GroupUser, 0)
	comIds = make([]int64, 0)
	DbEngin.Where("user_id = ?", userId).Find(&conconts)
	for _, v := range conconts {
		comIds = append(comIds, v.GroupId)
	}
	return comIds
}

// 获取单个群信息
func (server *ChatFriendsService) ShowCommunityID(dstId int64) (ret model.ChatGroup, err error) {
	com := model.ChatGroup{
		Id: dstId,
	}
	b, _ := DbEngin.Get(&com)
	if b == false {
		return ret, errors.New("该群不存在")
	} else {
		return com, nil
	}
}

// 批量获取群信息
func (service *ChatFriendsService) GetGroupInfo(ids []int64) []model.ChatGroup {
	tmp := make([]model.ChatGroup, 0)
	DbEngin.In("id", ids).Find(&tmp)
	return tmp
}

// 通过群id获取群成员id
func (service *ChatFriendsService) GetGroupUserByIds(groupId int64) (userIds []string) {
	conconts := make([]model.GroupUser, 0)
	userIds = make([]string, 0)
	DbEngin.Where("group_id = ?", groupId).Find(&conconts)
	fmt.Println(conconts, "------------")
	for _, v := range conconts {
		userIds = append(userIds, strconv.FormatInt(v.UserId, 10))
	}
	return userIds
}

//
func (s *ChatFriendsService) AddGroupChatDetailIds(arg []args.GroupChatDetailIds) error {
	for _, v := range arg {
		cot := model.ChatGroupDetailIds{
			UserId:  v.Userid,
			SendStatus: v.SendStatus,
			CreateTime: v.CreateTime,
		}
		_, err := DbEngin.InsertOne(&cot)
		if err != nil {
			return err
		}

	}
	return nil
}

// 好友备注
func (service *ChatFriendsService) FriendPetName(userid, friend_id int64,pet_name string) error {

	tmp := model.ChatFriends{}
	tmp.PetName = pet_name
	DbEngin.Where("user_id = ? and friend_id = ?", userid, friend_id).Cols("pet_name").Update(&tmp)

	return nil
}

// 修改群名称
func (s *ChatFriendsService) UpdateGroupName(group_id int64, nickname string) {
	tmp := model.ChatGroup{}
	tmp.Title = nickname
	DbEngin.ID(group_id).Cols("title").Update(&tmp)
}