package controllers

import (
	"fmt"
	"github.com/snowlyg/blog/relate"
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
* @api {get} /admin/docs/:id 根据id获取分类信息
* @apiName 根据id获取分类信息
* @apiGroup Docs
* @apiVersion 1.0.0
* @apiDescription 根据id获取分类信息
* @apiSampleRequest /admin/docs/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
 */
func GetDoc(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	relation := ctx.FormValue("relation")

	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			},
		},
		Relations: models.GetRelations(relation, relate.DocChapterRelate),
	}

	doc, err := models.GetDoc(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	_, _ = ctx.JSON(libs.ApiResource(200, docTransform(doc), "操作成功"))
}

/**
* @api {post} /admin/docs/ 新建分类
* @apiName 新建分类
* @apiGroup Docs
* @apiVersion 1.0.0
* @apiDescription 新建分类
* @apiSampleRequest /admin/docs/
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiDoc null
 */
func CreateDoc(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	doc := new(models.Doc)
	if err := ctx.ReadJSON(doc); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*doc)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	err = doc.CreateDoc()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if doc.ID == 0 {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, "操作失败"))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, docTransform(doc), "操作成功"))

}

/**
* @api {post} /admin/docs/:id/update 更新分类
* @apiName 更新分类
* @apiGroup Docs
* @apiVersion 1.0.0
* @apiDescription 更新分类
* @apiSampleRequest /admin/docs/:id/update
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiDoc null
 */
func UpdateDoc(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	aul := new(models.Doc)

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
	aul.ID = id
	err = models.UpdateDocById(id, aul)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error update doc: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if aul.ID == 0 {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, "操作失败"))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, docTransform(aul), "操作成功"))

}

/**
* @api {delete} /admin/docs/:id/delete 删除分类
* @apiName 删除分类
* @apiGroup Docs
* @apiVersion 1.0.0
* @apiDescription 删除分类
* @apiSampleRequest /admin/docs/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiDoc null
 */
func DeleteDoc(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteDocById(id)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /tts 获取所有的分类
* @apiName 获取所有的分类
* @apiGroup Docs
* @apiVersion 1.0.0
* @apiDescription 获取所有的分类
* @apiSampleRequest /tts
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiDoc null
 */
func GetAllDocs(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	s := GetCommonListSearch(ctx)

	docs, count, err := models.GetAllDocs(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	transform := docsTransform(docs)
	_, _ = ctx.JSON(libs.ApiResource(200, map[string]interface{}{"items": transform, "total": count, "limit": s.Limit}, "操作成功"))

}

func docsTransform(docs []*models.Doc) []*transformer.Doc {
	var rs []*transformer.Doc
	for _, doc := range docs {
		r := docTransform(doc)
		rs = append(rs, r)
	}
	return rs
}

func docTransform(doc *models.Doc) *transformer.Doc {
	r := &transformer.Doc{}
	g := gf.NewTransform(r, doc, time.RFC3339)
	_ = g.Transformer()
	if doc.Chapters != nil {
		transform := chaptersTransform(doc.Chapters)
		r.Chapters = transform
	}
	return r
}
