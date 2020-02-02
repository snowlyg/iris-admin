package controllers

import (
	"net/http"

	"IrisAdminApi/libs"
	"IrisAdminApi/models"
	"IrisAdminApi/validates"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

/**
* @api {post} /admin/login 用户登陆
* @apiName 用户登陆
* @apiGroup Admins
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
func AdminLogin(ctx iris.Context) {
	aul := new(validates.AdminLoginRequest)

	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil, e))
				return
			}
		}
	}

	ctx.Application().Logger().Infof("%s 登录系统", aul.Adminname)
	ctx.StatusCode(iris.StatusOK)

	admin := models.NewAdmin(0, aul.Adminname)
	admin.GetAdminByAdminname()

	response, status, msg := admin.CheckLogin(aul.Password)

	_, _ = ctx.JSON(ApiResource(status, response, msg))
	return

}

/**
* @api {get} /logout 用户退出登陆
* @apiName 用户退出登陆
* @apiGroup Admins
* @apiVersion 1.0.0
* @apiDescription 用户退出登陆
* @apiSampleRequest /logout
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func AdminLogout(ctx iris.Context) {
	aui := ctx.Values().GetString("auth_user_id")
	uid := uint(libs.ParseInt(aui, 0))
	models.AdminAdminLogout(uid)

	ctx.Application().Logger().Infof("%d 退出系统", uid)
	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "退出"))
}
