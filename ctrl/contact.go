package ctrl

import (
	"encoding/json"
	"net/http"
	"reptile-go/args"
	"reptile-go/model"
	"reptile-go/server"
	"reptile-go/util"
	"reptile-go/validates"
	"strconv"
)

var chatFriendsService server.ChatFriendsService
var contactValidate validates.ContactValidate

//添加好友
func Addfriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	contactValidates, err := contactValidate.ContactValidates(arg.Userid, arg.FriendId)
	if err != nil {
		util.RespFail(w, contactValidates)
		return
	}
	//调用service
	err = chatFriendsService.AddFriend(arg.Userid, arg.FriendId)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}
}
//删除好友双向删除
func Delfriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	contactValidates, err := contactValidate.ContactValidates(arg.Userid, arg.FriendId)
	if err != nil {
		util.RespFail(w, contactValidates)
		return
	}
	//调用service
	err = chatFriendsService.DelFriend(arg.Userid, arg.FriendId)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友删除成功")
	}
}

/**
@api {post} /contact/loadfriend 加载好友列表
*/

func LoadFriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	users := chatFriendsService.SearchFriend(arg.Userid)
	util.RespOkList(w, users, len(users))
}

/**
@api {post} /contact/createcommunity 创建群
@apiName 创建群
@apiGroup 群
@apiParam {String} name 群昵称
@apiParam {String} ownerid  用户ID
@apiParam {String} icon  群logo
@apiParam {String{可选}} icon  群logo

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"data": "",
	"msg": "xxx"
}
@apiError UserNotFound The id of the User was not found.

@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
*/
func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.ChatGroup
	util.Bind(r, &arg)
	if arg.AdminId == 0 || len(arg.Title) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	conn, err := chatFriendsService.CreateCommunity(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, conn, "创建群成功")
	}
}

// 加入群
/**
@api {post} /contact/joincommunity 加入群
@apiName 加入群
@apiGroup 群

@apiParam {Number} userid 用户ID
@apiParam {Number} dstid 群ID

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"data": "",
	"msg": "xxx"
}
@apiError UserNotFound The id of the User was not found.

@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
*/
//添加群成员
func JoinCommunity(w http.ResponseWriter, r *http.Request) {
	var arg []args.GroupArg
	s := r.Form.Get("data")
	err := json.Unmarshal([]byte(s), &arg)
	if err != nil {
		util.RespFail(w, "参数错误")
		return
	}

	err = chatFriendsService.JoinCommunity(arg)
	//todo 刷新用户的群组信息
	for _, v := range arg {
		AddGroupId(v.UserId, v.GroupId)
	}
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "success")
	}
}

/**
@api {post} /contact/loadcommunity 获取群列表
@apiName 获取群列表
@apiGroup 群

*/
func LoadCommunity(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	comunitys := chatFriendsService.SearchComunity(arg.Userid)
	util.RespOkList(w, comunitys, len(comunitys))
}

//获取单个群详情
func GetOneGroupInfo(w http.ResponseWriter, r *http.Request) {
	var arg args.GroupArg
	util.Bind(r, &arg)
	if arg.GroupId == 0 {
		util.RespFail(w, "参数错误")
		return
	}

	groupInfo, _ := chatFriendsService.ShowCommunityID(arg.GroupId)

	util.RespOk(w, groupInfo, "")
}

//批量获取群详情
func GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	//获取群ids
	comunitys := chatFriendsService.SearchComunityIds(arg.Userid)

	groupInfo := chatFriendsService.GetGroupInfo(comunitys)
	for k, v := range groupInfo {
		userIds := chatFriendsService.GetGroupUserByIds(v.Id)
		groupUsers := userService.GetUserByIds(userIds)
		groupInfo[k].GroupUsers = groupUsers
	}

	util.RespOkList(w, groupInfo, len(groupInfo))
}

func AddGroupChatDetailIds(w http.ResponseWriter, r *http.Request) {
	var arg []args.GroupChatDetailIds
	s := r.Form.Get("data")
	err := json.Unmarshal([]byte(s), &arg)
	if err != nil {
		util.RespFail(w, "参数错误")
		return
	}

	//调用service
	err = chatFriendsService.AddGroupChatDetailIds(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "")
	}
}

//好友备注
func FriendPetName(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	contactValidates, err := contactValidate.ContactValidates(arg.Userid, arg.FriendId)
	if err != nil {
		util.RespFail(w, contactValidates)
		return
	}
	//调用service
	err = chatFriendsService.FriendPetName(arg.Userid, arg.FriendId,arg.PetName)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "设置成功")
	}
}

//修改群名称
func UpdateGroupName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	group_id := r.PostForm.Get("group_id")
	nickname := r.PostForm.Get("nickname")
	if len(group_id) == 0 || len(nickname) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(group_id)
	chatFriendsService.UpdateGroupName(int64(id), nickname)
	util.RespOk(w, nil, "")
}