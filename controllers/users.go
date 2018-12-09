package controllers

import (
	"IrisApiProject/models"
	"fmt"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
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
func CGetProfile(ctx iris.Context) {
	aun := ctx.Values().Get("auth_user_name")
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, aun, "操作成功"))
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
func CGetUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	u := new(models.Users)
	u.ID = id
	user := u.GetUserById()

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, user, "操作成功"))
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
func CCreateUser(ctx iris.Context) {
	aul := new(models.UserJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.Type())
				fmt.Println(err.Param())
				fmt.Println()
			}
		} else {
			u := models.MCreateUser(aul)
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(apiResource(true, u, "操作成功"))
		}
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
func CUpdateUser(ctx iris.Context) {
	aul := new(models.UserJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.Type())
				fmt.Println(err.Param())
				fmt.Println()
			}
		} else {
			u := models.MUpdateUser(aul)
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(apiResource(true, u, "操作成功"))
		}
	}
}

/**
* @api {get} /admin/users/:id/frozen 冻结账号
* @apiName 冻结账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 冻结账号
* @apiSampleRequest /admin/users/:id/frozen
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CFrozenUser(ctx iris.Context) {

	id, _ := ctx.Params().GetUint("id")
	u := new(models.Users)
	u.ID = id

	is_frozen, msg := false, "冻结失败"
	if is_frozen = u.FrozenUserById(); is_frozen {
		msg = "冻结成功"
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(is_frozen, nil, msg))
}

/**
* @api {get} /admin/users/:id/refrozen 解冻用户
* @apiName 解冻用户
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 解冻用户
* @apiSampleRequest /admin/users/:id/refrozen
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CRefrozenUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	u := new(models.Users)
	u.ID = id

	is_frozen, msg := false, "解冻失败"
	if is_frozen = u.FrozenUserById(); is_frozen {
		msg = "解冻成功"
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(is_frozen, nil, msg))
}

/**
* @api {get} /admin/users/:id/aduit 设置负责人
* @apiName 设置负责人
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 设置负责人
* @apiSampleRequest /admin/users/:id/aduit
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CSetUserAudit(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	u := new(models.Users)
	u.ID = id
	u.SetAuditUserById()

	is_audit, msg := false, "设置失败"
	if is_audit = u.SetAuditUserById(); is_audit {
		msg = "设置成功"
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(is_audit, nil, msg))

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
func CDeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	u := new(models.Users)
	u.ID = id
	u.DeleteUserById()

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, nil, "删除成功"))

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
func CGetAllUsers(ctx iris.Context) {
	offset := t.ParseInt(ctx.FormValue("offset"), 1)
	limit := t.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	username := ctx.FormValue("username")
	orderBy := ctx.FormValue("orderBy")
	users := models.MGetAllUsers(name, username, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, users, "操作成功"))
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
func CGetAllClients(ctx iris.Context) {
	offset := t.ParseInt(ctx.FormValue("offset"), 1)
	limit := t.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	username := ctx.FormValue("username")
	orderBy := ctx.FormValue("orderBy")
	users := models.MGetAllClients(name, username, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, users, "操作成功"))
}
