package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/validates"
	"github.com/snowlyg/easygorm"
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
	ctx.StatusCode(iris.StatusOK)

	aul := new(validates.LoginRequest)

	if err := ctx.ReadJSON(aul); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	ctx.Application().Logger().Infof("%s 登录系统", aul.Username)

	s := &easygorm.Search{
		Fields: []*easygorm.Field{
			{
				Key:       "username",
				Condition: "=",
				Value:     aul.Username,
			},
		},
	}

	user, err := models.GetUser(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	response, code, msg := user.CheckLogin(aul.Password)

	_, _ = ctx.JSON(libs.ApiResource(code, response, msg))
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

	ctx.StatusCode(iris.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)
	conn := libs.GetRedisClusterClient()
	defer conn.Close()
	sess, err := models.GetRedisSessionV2(conn, value.Raw)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	if sess != nil {
		if err := sess.DelUserTokenCache(conn, value.Raw); err != nil {
			_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
			return
		}
	}

	ctx.Application().Logger().Infof("%d 退出系统", sess.UserId)
	_, _ = ctx.JSON(libs.ApiResource(200, nil, "退出"))
}

/**
* @api {get} /expire 刷新token
* @apiName 刷新token
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 刷新token
* @apiSampleRequest /expire
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UserExpire(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)
	conn := libs.GetRedisClusterClient()
	defer conn.Close()
	sess, err := models.GetRedisSessionV2(conn, value.Raw)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	if sess != nil {
		if err := sess.UpdateUserTokenCacheExpire(conn, value.Raw); err != nil {
			_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
			return
		}
	}

	_, _ = ctx.JSON(libs.ApiResource(200, nil, ""))
}
