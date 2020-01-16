package controllers

import (
	"net/http"

	"IrisAdminApi/models"
	"IrisAdminApi/tools"

	"github.com/kataras/iris/v12"
)

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
func UserLogin(ctx iris.Context) {
	aul := new(models.UserRequest)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(false, nil, "请求参数错误"))
	} else {
		if UserNameErr := validate.Var(aul.Username, "required,min=4,max=20"); UserNameErr != nil {
			ctx.StatusCode(iris.StatusOK)
			_, _ = ctx.JSON(ApiResource(false, nil, "用户名格式错误"))
			return
		} else if PwdErr := validate.Var(aul.Password, "required,min=5,max=20"); PwdErr != nil {
			ctx.StatusCode(iris.StatusOK)
			_, _ = ctx.JSON(ApiResource(false, nil, "密码格式错误"))
			return
		} else {
			ctx.Application().Logger().Infof("%s 登录系统", aul.Username)
			ctx.StatusCode(iris.StatusOK)
			response, status, msg := models.CheckLogin(aul.Username, aul.Password)
			_, _ = ctx.JSON(ApiResource(status, response, msg))
			return
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
func UserLogout(ctx iris.Context) {
	aui := ctx.Values().GetString("auth_user_id")
	uid := uint(tools.ParseInt(aui, 0))
	models.UserAdminLogout(uid)

	ctx.Application().Logger().Infof("%d 退出系统", uid)
	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "退出"))
}
