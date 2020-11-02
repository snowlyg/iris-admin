package controllers

import (
	"github.com/snowlyg/IrisAdminApi/libs"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/models"
	"github.com/snowlyg/IrisAdminApi/transformer"
	"github.com/snowlyg/IrisAdminApi/validates"
	gf "github.com/snowlyg/gotransformer"
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
	user, err := models.GetUserById(userId)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(libs.ApiResource(200, nil, err.Error()))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(libs.ApiResource(200, userTransform(user), ""))
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
	user, err := models.GetUserById(id)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(libs.ApiResource(200, nil, err.Error()))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(libs.ApiResource(200, userTransform(user), "操作成功"))
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
	ctx.StatusCode(iris.StatusOK)
	user := new(models.User)
	if err := ctx.ReadJSON(user); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*user)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	if err := user.CreateUser(); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
	}

	if user.ID == 0 {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, "操作失败"))
	} else {
		_, _ = ctx.JSON(libs.ApiResource(200, userTransform(user), "操作成功"))
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
	ctx.StatusCode(iris.StatusOK)
	user := new(models.User)
	if err := ctx.ReadJSON(user); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
	}

	err := validates.Validate.Struct(*user)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	user.ID = id
	if user.Username == "username" {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, "不能编辑管理员"))
		return
	}

	if err := user.UpdateUser(); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	_, _ = ctx.JSON(libs.ApiResource(200, userTransform(user), "操作成功"))

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
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	user, err := models.GetUserById(id)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	if user.Username == "username" {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, "不能删除管理员"))
		return
	}

	if err := models.DeleteUser(id); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	_, _ = ctx.JSON(libs.ApiResource(200, nil, "删除成功"))
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
	_, _ = ctx.JSON(libs.ApiResource(200, usersTransform(users), "操作成功"))
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

	roleIds := models.GetRolesForUser(user.ID)
	var ris []int
	var roles []*models.Role
	for _, roleId := range roleIds {
		ri, _ := strconv.Atoi(roleId)
		ris = append(ris, ri)
		role, err := models.GetRoleById(uint(ri))
		if err == nil {
			roles = append(roles, role)
		}

	}
	u.RoleIds = ris
	u.Roles = rolesTransform(roles)
	return u
}
