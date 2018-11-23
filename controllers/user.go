package controllers

import (
	"IrisYouQiKangApi/logic"
	"IrisYouQiKangApi/models"
	"github.com/kataras/iris"
	"net/http"
)

type AdminUserLogin struct {
	Username string
	Password string
}

/**
* @api {get} /admin/users/profile 获取登陆用户信息
* @apiName 获取登陆用户信息
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 获取登陆用户信息
* @apiSampleRequest /profile
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission 登陆用户
 */
func GetProfile(ctx iris.Context) {
	aun := ctx.Values().Get("auth_user_name")
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.ApiJson{State: true, Data: aun, Msg: "操作成功"})
}

/**
* @api {post} /admin/login 用户登陆
* @apiName 用户登陆
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 用户登陆
* @apiSampleRequest /login
* @apiParam {string} username 用户名
* @apiParam {string} password 密码
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UserAdminLogin(ctx iris.Context) {
	aul := new(AdminUserLogin)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(errorData(err))
	} else {
		err1 := validate.Var(aul.Username, "required,min=4,max=20")
		err2 := validate.Var(aul.Password, "required,min=5,max=20")
		if err1 != nil || err2 != nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(errorData(err1, err2))
		} else {
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(logic.UserAdminCheckLogin(aul.Username, aul.Password))
		}
	}
}

/**
* @api {get} /users 获取所有的账号
* @apiName 获取所有的账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 获取所有的账号
* @apiSampleRequest /users
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllUsers(ctx iris.Context) {
	cp := Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.GetAllUsers(kw, cp, mp))
}

/**
* @api {get} /clients 获取所有的客户联系人
* @apiName 获取所有的客户联系人
* @apiGroup Clients
* @apiVersion 1.0.0
* @apiDescription 获取所有的客户联系人
* @apiSampleRequest /clients
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllClients(ctx iris.Context) {
	cp := Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.GetAllClients(kw, cp, mp))
}

/**
* @api {get} /logout 用户退出登陆
* @apiName 用户退出登陆
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 用户退出登陆
* @apiSampleRequest /logout
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UserAdminLogout(ctx iris.Context) {
	json := models.ApiJson{}
	aui := ctx.Values().GetString("auth_user_id")

	uid := uint(Tools.ParseInt(aui, 0))

	json = logic.UserAdminLogout(uid)

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(json)
}
