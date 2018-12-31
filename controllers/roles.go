package controllers

import (
	"IrisApiProject/models"
	"IrisApiProject/tools"
	"fmt"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

/**
* @api {get} /admin/roles/:id 根据id获取角色信息
* @apiName 根据id获取角色信息
* @apiGroup Roles
* @apiVersion 1.0.0
* @apiDescription 根据id获取角色信息
* @apiSampleRequest /admin/roles/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetRole(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	role := models.GetRoleById(id)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(true, role, "操作成功"))
}

/**
* @api {post} /admin/roles/ 新建角色
* @apiName 新建角色
* @apiGroup Roles
* @apiVersion 1.0.0
* @apiDescription 新建角色
* @apiSampleRequest /admin/roles/
* @apiParam {string} name 角色名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CreateRole(ctx iris.Context) {

	aul := new(models.RoleJson)

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
			perms := []models.Permission{}
			u := models.CreateRole(aul, perms)
			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				ctx.JSON(ApiResource(false, u, "操作失败"))
			} else {
				ctx.JSON(ApiResource(true, u, "操作成功"))
			}
		}
	}
}

/**
* @api {post} /admin/roles/:id/update 更新角色
* @apiName 更新角色
* @apiGroup Roles
* @apiVersion 1.0.0
* @apiDescription 更新角色
* @apiSampleRequest /admin/roles/:id/update
* @apiParam {string} name 角色名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UpdateRole(ctx iris.Context) {
	aul := new(models.RoleJson)

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
			id, _ := ctx.Params().GetInt("id")
			uid := uint(id)

			u := models.UpdateRole(aul, uid)
			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				ctx.JSON(ApiResource(false, u, "操作失败"))
			} else {
				ctx.JSON(ApiResource(true, u, "操作成功"))
			}
		}
	}
}

/**
* @api {delete} /admin/roles/:id/delete 删除角色
* @apiName 删除角色
* @apiGroup Roles
* @apiVersion 1.0.0
* @apiDescription 删除角色
* @apiSampleRequest /admin/roles/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func DeleteRole(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	models.DeleteRoleById(id)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(true, nil, "删除成功"))
}

/**
* @api {get} /roles 获取所有的角色
* @apiName 获取所有的角色
* @apiGroup Roles
* @apiVersion 1.0.0
* @apiDescription 获取所有的角色
* @apiSampleRequest /roles
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllRoles(ctx iris.Context) {
	offset := tools.Tool.ParseInt(ctx.FormValue("offset"), 1)
	limit := tools.Tool.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	roles := models.GetAllRoles(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(true, roles, "操作成功"))
}
