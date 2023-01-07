package ctrl

import (
	"math/rand"
	"net/http"
	"reptile-go/model"
	"reptile-go/server"
	"reptile-go/util"
	"strconv"
	"strings"
	"time"
)

var userService server.UserService

// 登录
/**
@api {post} /index/user/login 登录
@apiName Login
@apiGroup 登录
@apiParam {String} [mobile='123546'] mobile 账号
@apiParam {String} passwd 密码
@apiHeaderExample {json} Header-Example:
	{
		"Content-Type":"application/x-www-form-urlencoded"
	}
@apiSuccess {Number} code 状态码.
@apiSuccess {Object} data  json数据.
@apiSuccess {String} msg  提示.
@apiError  {Number} code 状态码.
@apiError  {String} msg  提示.
*/

func UserLogin(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	if r.Method == http.MethodPost {
		r.ParseForm()
		//1.获取前端传递过来的参数
		username := r.PostForm.Get("username")
		plainpwd := r.PostForm.Get("password")
		if len(username) == 0 || len(plainpwd) == 0 {
			util.RespFail(w, "参数错误")
			return
		}
		user, err := userService.Login(username, plainpwd)
		if err != nil {
			util.RespFail(w, err.Error())
		} else {
			// token
			tokenString, _ := util.GenToken(user.Username)
			util.RespOk(w, user, tokenString)
		}
	}
}

// 注册
/**
@api {post} /index/user/register 注册
@apiGroup 注册
@apiParam {String} [mobile='123546'] mobile 账号
@apiParam {String} passwd 密码
@apiParam {String} uuid key
@apiParam {String{5}} code 验证码
@apiHeaderExample {json} Header-Example:
	{
		"Content-Type":"application/x-www-form-urlencoded"
	}
@apiSuccess {Number} code 状态码.
@apiSuccess {Object} data  json数据.
@apiSuccess {String} msg  提示.
@apiError  {Number} code 状态码.
@apiError  {String} msg  提示.
*/
func UserRegister(w http.ResponseWriter, r *http.Request) {
	//1.获取前端传递过来的参数
	// 解析参数
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.PostForm.Get("username")
		nickname := r.PostForm.Get("nickname")
		plainpwd := r.PostForm.Get("passwd")
		uuid := r.PostForm.Get("uuid")
		code := r.PostForm.Get("code")
		if len(username) == 0 || len(plainpwd) == 0 || len(code) == 0 || len(uuid) == 0 {
			util.RespFail(w, "参数错误")
			return
		}
		// 检验验证码
		err := util.CaptchaVerifyHandle(uuid, code)
		if err != nil {
			util.RespFail(w, err.Error())
			return
		}
		rand.Seed(time.Now().UnixNano()) // 设置种子数为当前时间
		//nickname := fmt.Sprintf("user%06d", rand.Int31())
		avatar := ""
		sex := model.SEX_UNKNOW
		user, err := userService.Register(username, plainpwd, nickname, avatar, sex)
		if err != nil {
			util.RespFail(w, err.Error())
		} else {
			util.RespOk(w, user, "")
		}
	}
}

//修改用户数据
/**
@api {post} /user/updateUser 修改用户数据
@apiName GetUser
@apiGroup 用户
@apiParam {Number} userid Users unique ID.
@apiSuccess {String} avatar 头像.
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
func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.PostForm.Get("userid")
	avatar := r.PostForm.Get("avatar")
	if len(userid) == 0 || len(avatar) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(userid)
	userService.UserInfo(int64(id), avatar)
	util.RespOk(w, nil, "")
}

//更新用户昵称
func UpdateUserNickname(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.PostForm.Get("userid")
	nickname := r.PostForm.Get("nickname")
	if len(userid) == 0 || len(nickname) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(userid)
	userService.UserNickname(int64(id), nickname)
	util.RespOk(w, nil, "")
}

//获取用户数据
func GetUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.PostForm.Get("userid")
	if len(userid) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(userid)
	user, _ := userService.GetUser(int64(id))
	util.RespOk(w, user, "")
}

//通过手机账号获取用户数据
func GetUserByName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobile := r.PostForm.Get("mobile")
	if len(mobile) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(mobile)
	user, _ := userService.GetUserByName(int64(id))
	util.RespOk(w, user, "")
}

//批量获取用户数据
func GetUserByIds(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := r.PostForm.Get("ids")
	ids_arr := strings.Split(ids, ",")
	user := userService.GetUserByIds(ids_arr)
	util.RespOkList(w, user, len(user))
}

//更新用户密码
func Editpwd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.PostForm.Get("userId")
	password := r.PostForm.Get("newPassword")
	if len(userid) == 0 || len(password) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(userid)
	userService.EditUserPwd(int64(id), password)
	util.RespOk(w, nil, "成功")
}

//退出登陆
func UpdateUserOnline(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.PostForm.Get("userid")
	online := 0
	if len(userid) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	id, _ := strconv.Atoi(userid)
	userService.UserOnline(int64(id), online)
	util.RespOk(w, nil, "")
}

/**
 * @api {get} /index/getCaptcha 获取验证码
 * @apiName registered
 * @apiGroup 注册
 * @apiSuccess {Number} code 状态码.
 * @apiSuccess {String} data  base64图片字符串.
 * @apiSuccess {String} id  字符串Key.
 * @apiSuccess {String} msg  提示.
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "code": 0,
 *       "data": "xxxxxx",
 *       "id": "xxxxxx",
 *       "msg": "xxxxx",
 *     }
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 404 Not Found
 *     {
 *       "code": -1,
 *       "msg": "xxxxx",
 *     }
 */

func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	util.GenerateCaptchaHandler(w, r)
}
