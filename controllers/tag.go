package controllers

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/libs"
	"github.com/snowlyg/IrisAdminApi/models"
	"github.com/snowlyg/IrisAdminApi/transformer"
	"github.com/snowlyg/IrisAdminApi/validates"
	gf "github.com/snowlyg/gotransformer"
)

/**
* @api {get} /admin/tags/:id 根据id获取权限信息
* @apiName 根据id获取权限信息
* @apiGroup Tags
* @apiVersion 1.0.0
* @apiDescription 根据id获取权限信息
* @apiSampleRequest /admin/tags/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
 */
func GetTag(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	tag, err := models.GetTagById(id)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, tagTransform(tag), "操作成功"))
}

/**
* @api {post} /admin/tags/ 新建权限
* @apiName 新建权限
* @apiGroup Tags
* @apiVersion 1.0.0
* @apiDescription 新建权限
* @apiSampleRequest /admin/tags/
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiTag null
 */
func CreateTag(ctx iris.Context) {
	tag := new(models.Tag)
	if err := ctx.ReadJSON(tag); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*tag)
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

	err = tag.CreateTag()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if tag.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, tag, "操作失败"))
	} else {
		_, _ = ctx.JSON(ApiResource(200, tagTransform(tag), "操作成功"))
	}

}

/**
* @api {post} /admin/tags/:id/update 更新权限
* @apiName 更新权限
* @apiGroup Tags
* @apiVersion 1.0.0
* @apiDescription 更新权限
* @apiSampleRequest /admin/tags/:id/update
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiTag null
 */
func UpdateTag(ctx iris.Context) {
	aul := new(models.Tag)

	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*aul)
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
	err = models.UpdateTagById(id, aul)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if aul.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, nil, "操作失败"))
	} else {
		_, _ = ctx.JSON(ApiResource(200, tagTransform(aul), "操作成功"))
	}

}

/**
* @api {delete} /admin/tags/:id/delete 删除权限
* @apiName 删除权限
* @apiGroup Tags
* @apiVersion 1.0.0
* @apiDescription 删除权限
* @apiSampleRequest /admin/tags/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiTag null
 */
func DeleteTag(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteTagById(id)
	if err != nil {

		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /tags 获取所有的权限
* @apiName 获取所有的权限
* @apiGroup Tags
* @apiVersion 1.0.0
* @apiDescription 获取所有的权限
* @apiSampleRequest /tags
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiTag null
 */
func GetAllTags(ctx iris.Context) {
	offset := libs.ParseInt(ctx.URLParam("offset"), 1)
	limit := libs.ParseInt(ctx.URLParam("limit"), 20)
	name := ctx.FormValue("searchStr")
	orderBy := ctx.FormValue("orderBy")

	tags, err := models.GetAllTags(name, orderBy, offset, limit)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, tagsTransform(tags), "操作成功"))
}

func tagsTransform(tags []*models.Tag) []*transformer.Tag {
	var rs []*transformer.Tag
	for _, tag := range tags {
		r := tagTransform(tag)
		rs = append(rs, r)
	}
	return rs
}

func tagTransform(tag *models.Tag) *transformer.Tag {
	r := &transformer.Tag{}
	g := gf.NewTransform(r, tag, time.RFC3339)
	_ = g.Transformer()
	return r
}
