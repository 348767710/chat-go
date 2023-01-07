package ctrl

import (
	"encoding/json"
	"net/http"
	"reptile-go/args"
	"reptile-go/model"
	"reptile-go/server"
	"reptile-go/util"
	"strconv"
	"strings"
)

var messageService server.MessageService

// 获取消息记录
/**
@api {post} /message/chathistory 获取消息记录
@apiName 消息记录
@apiGroup 聊天消息

@apiParam {Number} userid 用户ID
@apiParam {Number} dstid 好友ID/群ID
@apiParam {Number} cmd 单/群聊

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"data": "",
	"msg": "xxx"
}

@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
@apiUse CommonError
*/
func ChatHistory(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 || arg.Dstid == 0 || arg.Type == "" {
		util.RespFail(w, "参数错误")
		return
	}
	var chat string
	if arg.Type == model.CMD_ROOM_MSG {
		chat = "chat_11"
	} else {
		chat = "chat_10"
	}
	lRange, _ := server.LRange(chat, int64(arg.GetPageFrom()), int64(arg.GetPageSize()))
	var aa = false
	if aa {
		var msg model.Message
		msgList := make([]model.Message, 0)
		for _, data := range lRange {
			json.Unmarshal([]byte(data), &msg)
			if arg.Type == model.CMD_ROOM_MSG {
				if msg.Uid != 0 && arg.Dstid == msg.Dstid {
					//msg.Createat = time.Now().Unix()
					//msgList = append(msgList, msg)
					msgList = append([]model.Message{msg}, msgList...)
				}
			} else {
				//(userid = ? and dstid = ?) or (dstid = ? and userid = ?)
				if msg.Uid != 0 && (msg.Dstid == arg.Userid && msg.Uid == arg.Dstid) || (arg.Userid == msg.Uid && arg.Dstid == msg.Dstid) {
					//msg.Createat = time.Now().Unix()
					//msgList = append(msgList, msg)
					msgList = append([]model.Message{msg}, msgList...)
				}
			}
		}
		util.RespOkList(w, msgList, len(msgList))
		return
	} else {
		history := messageService.GetChatHistory(arg.Userid, arg.Dstid, arg.Type, arg.GetPageFrom(), arg.GetPageSize())
		util.RespOkList(w, history, len(history))
	}
}

// 获取用户未阅读的历史消息
func GetNoReadHistory(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	history := messageService.GetUserNoReadHistory(arg.Userid)
	util.RespOkList(w, history, len(history))
}

// 更新用户历史消息已读回执
func UpdateNoReadHistory(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	ids := r.PostForm.Get("ids")
	ids_arr := strings.Split(ids, ",")
	messageService.UpdateNoReadHistory(ids_arr)
	util.RespOk(w, nil, "")
}

// 获取用户未阅读的群组历史消息
func GetNoReadGroupHistory(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	groups := messageService.GetNoReadGroupHistory(arg.Userid)
	var msg model.ChatGroupDetail
	msgList := make([]string, 0)
	historyList := make([]model.ChatGroupDetail, 0)
	//resList := [...][]string{}

	for _, v := range groups {
		msgList, _ = server.LRange("chat_group_redis_key_"+strconv.FormatInt(v.GroupId, 10), -100, -1)
		for _, data := range msgList {
			json.Unmarshal([]byte(data), &msg)
			if msg.ToId != 0 {
				historyList = append(historyList, msg)
			}
		}
		//resList = append(resList,historyList)
	}

	util.RespOkList(w, historyList, len(historyList))
}

// 更新用户群组历史消息已读回执
func UpdateNoReadGroupHistory(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	ids := r.PostForm.Get("ids")
	ids_arr := strings.Split(ids, ",")
	messageService.UpdateNoReadGroupHistory(ids_arr)
	util.RespOk(w, nil, "")
}
