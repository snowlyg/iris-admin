package controllers

import (
	"IrisYouQiKangApi/logic"
	"IrisYouQiKangApi/models"
	"github.com/kataras/iris"
	"net/http"
)

var apiJson models.ApiJson

/**
* @api {post} /admin/login 用户登陆
* @apiName 用户登陆
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 用户登陆
* @apiSampleRequest /admin/login
* @apiParam {string} username 用户名
* @apiParam {string} password 密码
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UserAdminLogin(ctx iris.Context) {
	aul := new(models.AdminUserLogin)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		apiJson.Msg = "请求参数错误"
		ctx.JSON(apiJson)
	} else {
		if UserNameErr := validate.Var(aul.Username, "required,min=4,max=20"); UserNameErr != nil {
			ctx.StatusCode(iris.StatusOK)
			apiJson.Msg = "用户名格式错误"
			ctx.JSON(apiJson)
		} else if PwdErr := validate.Var(aul.Password, "required,min=5,max=20"); PwdErr != nil {
			ctx.StatusCode(iris.StatusOK)
			apiJson.Msg = "密码格式错误"
			ctx.JSON(apiJson)
		} else {
			ctx.StatusCode(iris.StatusOK)
			response, status, msg := logic.UserAdminCheckLogin(aul.Username, aul.Password)
			apiJson.Msg = msg
			apiJson.Status = status
			apiJson.Data = response
			ctx.JSON(apiJson)
		}
	}
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
