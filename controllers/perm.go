package controllers

import (
	"fmt"
	"github.com/snowlyg/blog/libs/easygorm"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/transformer"
	"github.com/snowlyg/blog/validates"
	gf "github.com/snowlyg/gotransformer"
)

/**
* @api {get} /admin/permissions/:id 根据id获取权限信息
* @apiName 根据id获取权限信息
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 根据id获取权限信息
* @apiSampleRequest /admin/permissions/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetPermission(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	s := &easygorm.Search{
		Fields: []*easygorm.Field{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			},
		},
	}
	perm, err := models.GetPermission(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	_, _ = ctx.JSON(libs.ApiResource(200, permTransform(perm), "操作成功"))
}

/**
* @api {post} /admin/permissions/ 新建权限
* @apiName 新建权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 新建权限
* @apiSampleRequest /admin/permissions/
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CreatePermission(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	perm := new(models.Permission)
	if err := ctx.ReadJSON(perm); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*perm)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	err = perm.CreatePermission()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	if perm.ID == 0 {
		_, _ = ctx.JSON(libs.ApiResource(400, perm, "操作失败"))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, permTransform(perm), "操作成功"))

}

/**
* @api {post} /admin/permissions/:id/update 更新权限
* @apiName 更新权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 更新权限
* @apiSampleRequest /admin/permissions/:id/update
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UpdatePermission(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	aul := new(models.Permission)

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

	id, _ := ctx.Params().GetUint("id")
	err = models.UpdatePermission(id, aul)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error update prem: %s", err.Error())))
		return
	}

	_, _ = ctx.JSON(libs.ApiResource(200, permTransform(aul), "操作成功"))

}

/**
* @api {delete} /admin/permissions/:id/delete 删除权限
* @apiName 删除权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 删除权限
* @apiSampleRequest /admin/permissions/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func DeletePermission(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	err := models.DeletePermissionById(id)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /permissions 获取所有的权限
* @apiName 获取所有的权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 获取所有的权限
* @apiSampleRequest /permissions
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllPermissions(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	s := GetCommonListSearch(ctx)
	permissions, count, err := models.GetAllPermissions(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	transform := permsTransform(permissions)
	_, _ = ctx.JSON(libs.ApiResource(200, map[string]interface{}{"items": transform, "total": count, "limit": s.Limit}, "操作成功"))

}

func permsTransform(perms []*models.Permission) []*transformer.Permission {
	var rs []*transformer.Permission
	for _, perm := range perms {
		r := permTransform(perm)
		rs = append(rs, r)
	}
	return rs
}

func permTransform(perm *models.Permission) *transformer.Permission {
	r := &transformer.Permission{}
	g := gf.NewTransform(r, perm, time.RFC3339)
	_ = g.Transformer()
	return r
}
