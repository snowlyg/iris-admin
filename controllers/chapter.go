package controllers

import (
	"fmt"
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
* @api {get} /chapters/:id 根据id获取文章信息
* @apiName 根据id获取文章信息
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 根据id获取文章信息
* @apiSampleRequest /chapters/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetPublishedChapter(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	relation := ctx.FormValue("relation")
	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			}, {
				Key:       "status",
				Condition: "=",
				Value:     "published",
			},
		},
		Relations: models.GetRelations(relation, nil),
	}
	chapter, err := models.GetChapter(s)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err = chapter.ReadChapter(ctx.Request())
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	rr := chapterTransform(chapter)
	_, _ = ctx.JSON(ApiResource(200, rr, "操作成功"))
}

/**
* @api {get} /admin/chapters/:id 根据id获取分类信息
* @apiName 根据id获取分类信息
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 根据id获取分类信息
* @apiSampleRequest /admin/chapters/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
 */
func GetChapter(ctx iris.Context) {

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
		Relations: models.GetRelations(relation, nil),
	}
	chapter, err := models.GetChapter(s)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}

	_, _ = ctx.JSON(ApiResource(200, chapterTransform(chapter), "操作成功"))
}

/**
* @api {post} /admin/chapters/ 新建分类
* @apiName 新建分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 新建分类
* @apiSampleRequest /admin/chapters/
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func CreateChapter(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	chapter := new(models.Chapter)
	if err := ctx.ReadJSON(chapter); err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(chapter)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(ApiResource(400, nil, e))
				return
			}
		}
	}

	err = chapter.CreateChapter()
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	if chapter.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, nil, "操作失败"))
	}
	_, _ = ctx.JSON(ApiResource(200, chapterTransform(chapter), "操作成功"))

}

/**
* @api {post} /admin/chapters/:id/update 更新分类
* @apiName 更新分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 更新分类
* @apiSampleRequest /admin/chapters/:id/update
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func UpdateChapter(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	aul := new(models.Chapter)

	if err := ctx.ReadJSON(aul); err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(ApiResource(400, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	aul.ID = id
	err = models.UpdateChapterById(id, aul)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error update chapter: %s", err.Error())))
		return
	}
	if aul.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, nil, "操作失败"))
	}
	_, _ = ctx.JSON(ApiResource(200, chapterTransform(aul), "操作成功"))

}

/**
* @api {delete} /admin/chapters/:id/delete 删除分类
* @apiName 删除分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 删除分类
* @apiSampleRequest /admin/chapters/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func DeleteChapter(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteChapterById(id)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}
	_, _ = ctx.JSON(ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /tts 获取所有的分类
* @apiName 获取所有的分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 获取所有的分类
* @apiSampleRequest /tts
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func GetAllChapters(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	offset := libs.ParseInt(ctx.URLParam("page"), 1)
	limit := libs.ParseInt(ctx.URLParam("limit"), 20)
	docId := libs.ParseInt(ctx.URLParam("docId"), 0)
	orderBy := ctx.FormValue("orderBy")
	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "doc_id",
				Condition: "=",
				Value:     docId,
			},
		},
		Offset:  offset,
		Limit:   limit,
		OrderBy: orderBy,
	}
	chapters, count, err := models.GetAllChapters(s)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}

	transform := chaptersTransform(chapters)
	_, _ = ctx.JSON(ApiResource(200, map[string]interface{}{"items": transform, "total": count, "limit": limit}, "操作成功"))
}

/**
* @api {get} /chapter/like/:id 根据id点赞文章
* @apiName 根据id点赞文章
* @apiGroup Chapter
* @apiVersion 1.0.0
* @apiDescription 根据id点赞文章
* @apiSampleRequest /chapter/like/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetPublishedChapterLike(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			}, {
				Key:       "status",
				Condition: "=",
				Value:     "published",
			},
		},
	}
	chapter, err := models.GetChapter(s)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err = chapter.LikeChapter()
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	rr := chapterTransform(chapter)
	_, _ = ctx.JSON(ApiResource(200, rr, "操作成功"))
}

/**
* @api {get} /chapter 获取所有的文章
* @apiName 获取所有的文章
* @apiGroup Chapter
* @apiVersion 1.0.0
* @apiDescription 获取所有的文章
* @apiSampleRequest /chapter
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllPublishedChapters(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	offset := libs.ParseInt(ctx.FormValue("page"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	docId := libs.ParseInt(ctx.FormValue("docId"), 0)
	orderBy := ctx.FormValue("orderBy")
	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "doc_id",
				Condition: "=",
				Value:     docId,
			}, {
				Key:       "status",
				Condition: "=",
				Value:     "published",
			},
		},
		Offset:  offset,
		Limit:   limit,
		OrderBy: orderBy,
	}
	chapters, count, err := models.GetAllChapters(s)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	transform := chaptersTransform(chapters)
	_, _ = ctx.JSON(ApiResource(200, map[string]interface{}{"items": transform, "total": count, "limit": limit}, "操作成功"))
}

func chaptersTransform(chapters []*models.Chapter) []*transformer.Chapter {
	var rs []*transformer.Chapter
	for _, chapter := range chapters {
		r := chapterTransform(chapter)
		rs = append(rs, r)
	}
	return rs
}

func chapterTransform(chapter *models.Chapter) *transformer.Chapter {
	r := &transformer.Chapter{}
	g := gf.NewTransform(r, chapter, time.RFC3339)
	_ = g.Transformer()
	if chapter.Doc != nil {
		transform := docTransform(chapter.Doc)
		r.Doc = *transform
	}
	return r
}
