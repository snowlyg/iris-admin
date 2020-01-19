package controllers

import (
	"strconv"
	"time"

	"IrisAdminApi/models"
	"IrisAdminApi/transformer"
	"github.com/go-playground/validator/v10"
	gf "github.com/snowlyg/gotransformer"

	"github.com/kataras/iris/v12"
)

/**
* @api {get} /admin/users/profile 获取登陆用户信息
* @apiName 获取登陆用户信息
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 获取登陆用户信息
* @apiSampleRequest /admin/users/profile
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission 登陆用户
 */
func GetProfile(ctx iris.Context) {
	userId := ctx.Values().Get("auth_user_id").(uint)
	user := models.GetUserById(userId)
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, userTransform(user), "操作成功"))
}

/**
* @api {get} /admin/users/:id 根据id获取用户信息
* @apiName 根据id获取用户信息
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 根据id获取用户信息
* @apiSampleRequest /admin/users/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission 登陆用户
 */
func GetUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	user := models.GetUserById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, userTransform(user), "操作成功"))
}

/**
* @api {post} /admin/users/ 新建账号
* @apiName 新建账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 新建账号
* @apiSampleRequest /admin/users/
* @apiParam {string} username 用户名
* @apiParam {string} password 密码
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CreateUser(ctx iris.Context) {

	aul := new(models.UserRequest)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	err := Validate.Struct(models.UserRequest{})
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil, e))
				return
			}
		}
	}

	u := models.CreateUser(aul)
	ctx.StatusCode(iris.StatusOK)
	if u.ID == 0 {
		_, _ = ctx.JSON(ApiResource(false, u, "操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ApiResource(true, nil, "操作成功"))
		return
	}

}

/**
* @api {post} /admin/users/:id/update 更新账号
* @apiName 更新账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 更新账号
* @apiSampleRequest /admin/users/:id/update
* @apiParam {string} username 用户名
* @apiParam {string} password 密码
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UpdateUser(ctx iris.Context) {
	aul := new(models.UserRequest)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
	}

	err := Validate.Struct(models.UserRequest{})
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	user := models.GetUserById(id)
	if user.Username == "username" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, "不能编辑管理员"))
		return
	}

	u := models.UpdateUser(aul, id)
	ctx.StatusCode(iris.StatusOK)
	if u.ID == 0 {
		_, _ = ctx.JSON(ApiResource(false, u, "操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ApiResource(true, nil, "操作成功"))
		return
	}

}

/**
* @api {delete} /admin/users/:id/delete 删除用户
* @apiName 删除用户
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 删除用户
* @apiSampleRequest /admin/users/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func DeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	user := models.GetUserById(id)
	if user.Username == "username" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, "不能删除管理员"))
		return
	}

	models.DeleteUserById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "删除成功"))
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
	offset := ctx.URLParamIntDefault("offset", 1)
	limit := ctx.URLParamIntDefault("limit", 15)
	name := ctx.URLParam("name")
	orderBy := ctx.URLParam("orderBy")

	users := models.GetAllUsers(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, usersTransform(users), "操作成功"))
}

func usersTransform(users []*models.User) []*transformer.User {
	var us []*transformer.User
	for _, user := range users {
		u := userTransform(user)
		us = append(us, u)
	}
	return us
}

func userTransform(user *models.User) *transformer.User {
	u := &transformer.User{}
	g := gf.NewTransform(u, user, time.RFC3339)
	_ = g.Transformer()

	roleIds, _ := models.Enforcer.GetRolesForUser(strconv.FormatUint(uint64(user.ID), 10))
	var ris []int
	for _, roleId := range roleIds {
		ri, _ := strconv.Atoi(roleId)
		ris = append(ris, ri)
	}
	u.RoleIds = ris
	return u
}
