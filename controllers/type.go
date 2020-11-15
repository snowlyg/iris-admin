package controllers

import (
	"fmt"
	"github.com/snowlyg/easygorm"
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
* @api {get} /admin/types/:id 根据id获取分类信息
* @apiName 根据id获取分类信息
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 根据id获取分类信息
* @apiSampleRequest /admin/types/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
 */
func GetType(ctx iris.Context) {
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
	tt, err := models.GetType(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	_, _ = ctx.JSON(libs.ApiResource(200, ttTransform(tt), "操作成功"))
}

/**
* @api {post} /admin/types/ 新建分类
* @apiName 新建分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 新建分类
* @apiSampleRequest /admin/types/
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func CreateType(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	tt := new(models.Type)
	if err := ctx.ReadJSON(tt); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*tt)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	err = tt.CreateType()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error create type: %s", err.Error())))
		return
	}

	if tt.ID == 0 {
		_, _ = ctx.JSON(libs.ApiResource(400, tt, "操作失败"))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, ttTransform(tt), "操作成功"))

}

/**
* @api {post} /admin/types/:id/update 更新分类
* @apiName 更新分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 更新分类
* @apiSampleRequest /admin/types/:id/update
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func UpdateType(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	aul := new(models.Type)

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
	err = models.UpdateTypeById(id, aul)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error update type: %s", err.Error())))
		return
	}

	_, _ = ctx.JSON(libs.ApiResource(200, ttTransform(aul), "操作成功"))

}

/**
* @api {delete} /admin/types/:id/delete 删除分类
* @apiName 删除分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 删除分类
* @apiSampleRequest /admin/types/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func DeleteType(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteTypeById(id)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /types 获取所有的分类
* @apiName 获取所有的分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 获取所有的分类
* @apiSampleRequest /types
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func GetAllTypes(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	s := GetCommonListSearch(ctx)
	tts, count, err := models.GetAllTypes(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	transform := ttsTransform(tts)
	_, _ = ctx.JSON(libs.ApiResource(200, map[string]interface{}{"items": transform, "total": count, "limit": s.Limit}, "操作成功"))

}

func ttsTransform(tts []*models.Type) []*transformer.Type {
	var rs []*transformer.Type
	for _, tt := range tts {
		r := ttTransform(tt)
		rs = append(rs, r)
	}
	return rs
}

func ttTransform(tt *models.Type) *transformer.Type {
	r := &transformer.Type{}
	g := gf.NewTransform(r, tt, time.RFC3339)
	_ = g.Transformer()
	return r
}
