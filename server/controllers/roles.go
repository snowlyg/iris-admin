package controllers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/server/libs"
	"github.com/snowlyg/IrisAdminApi/server/models"
	"github.com/snowlyg/IrisAdminApi/server/transformer"
	"github.com/snowlyg/IrisAdminApi/server/validates"
	gf "github.com/snowlyg/gotransformer"
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
	role := models.NewRole(id, "")
	role.GetRoleById()

	ctx.StatusCode(iris.StatusOK)

	rr := roleTransform(role)
	rr.Perms = permsTransform(role.RolePermisions())
	_, _ = ctx.JSON(ApiResource(200, rr, "操作成功"))
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
	role := new(models.Role)

	if err := ctx.ReadJSON(role); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*role)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(400, nil, e))
				return
			}
		}
	}

	err = role.CreateRole()
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, role, err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	if role.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, role, "操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ApiResource(200, role, "操作成功"))
		return
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
	role := new(models.Role)
	if err := ctx.ReadJSON(role); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*role)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(400, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	role.ID = id
	if role.Name == "admin" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, "不能编辑管理员角色"))
		return
	}

	role.UpdateRole(role)
	ctx.StatusCode(iris.StatusOK)
	if role.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, role, "操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ApiResource(200, nil, "操作成功"))
		return
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
	role := models.NewRole(id, "")
	role.GetRoleById()
	if role.Name == "admin" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(200, nil, "不能删除管理员角色"))
		return
	}

	role.DeleteRoleById()

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, nil, "删除成功"))
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
	offset := libs.ParseInt(ctx.FormValue("offset"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	roles, err := models.GetAllRoles(name, orderBy, offset, limit)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, rolesTransform(roles), "操作成功"))
}

func rolesTransform(roles []*models.Role) []*transformer.Role {
	var rs []*transformer.Role
	for _, role := range roles {
		r := roleTransform(role)
		rs = append(rs, r)
	}
	return rs
}

func roleTransform(role *models.Role) *transformer.Role {
	r := &transformer.Role{}
	g := gf.NewTransform(r, role, time.RFC3339)
	_ = g.Transformer()
	r.Perms = permsTransform(role.RolePermisions())
	return r
}
