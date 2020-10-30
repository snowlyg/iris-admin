package controllers

import (
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
* @api {get} /articles/:id 根据id获取文章信息
* @apiName 根据id获取文章信息
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 根据id获取文章信息
* @apiSampleRequest /articles/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetPublishedArticle(ctx iris.Context) {
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
	article, err := models.GetArticle(s)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err = article.ReadArticle(ctx.Request())
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	rr := articleTransform(article)
	_, _ = ctx.JSON(ApiResource(200, rr, "操作成功"))
}

/**
* @api {get} /articles/like/:id 根据id点赞文章
* @apiName 根据id点赞文章
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 根据id点赞文章
* @apiSampleRequest /articles/like/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetPublishedArticleLike(ctx iris.Context) {
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
	article, err := models.GetArticle(s)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err = article.LikeArticle()
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	rr := articleTransform(article)
	_, _ = ctx.JSON(ApiResource(200, rr, "操作成功"))
}

/**
* @api {get} /admin/articles/:id 根据id获取文章信息
* @apiName 根据id获取文章信息
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 根据id获取文章信息
* @apiSampleRequest /articles/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetArticle(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	relation := ctx.FormValue("relation")
	search := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "id",
				Value:     id,
				Condition: "=",
			},
		},
		Relations: models.GetRelations(relation, nil),
	}
	article, err := models.GetArticle(search)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	rr := articleTransform(article)
	_, _ = ctx.JSON(ApiResource(200, rr, "操作成功"))
}

/**
* @api {post} /admin/articles/ 新建文章
* @apiName 新建文章
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 新建文章
* @apiSampleRequest /admin/articles/
* @apiParam {string} name 文章名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CreateArticle(ctx iris.Context) {
	article := new(models.Article)

	if err := ctx.ReadJSON(article); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(article)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(200, nil, e))
				return
			}
		}
	}

	err = article.CreateArticle()
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	if article.ID == 0 {
		_, _ = ctx.JSON(ApiResource(200, nil, "操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ApiResource(200, articleTransform(article), "操作成功"))
		return
	}

}

/**
* @api {post} /admin/articles/:id/update 更新文章
* @apiName 更新文章
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 更新文章
* @apiSampleRequest /admin/articles/:id/update
* @apiParam {string} name 文章名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UpdateArticle(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)

	article := new(models.Article)
	if err := ctx.ReadJSON(article); err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(article)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(ApiResource(200, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	err = models.UpdateArticle(id, article)

	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	_, _ = ctx.JSON(ApiResource(200, articleTransform(article), "操作成功"))
	return

}

/**
* @api {delete} /admin/articles/:id/delete 删除文章
* @apiName 删除文章
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 删除文章
* @apiSampleRequest /admin/articles/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func DeleteArticle(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteArticleById(id)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /articles 获取所有的文章
* @apiName 获取所有的文章
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 获取所有的文章
* @apiSampleRequest /articles
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllPublishedArticles(ctx iris.Context) {
	offset := libs.ParseInt(ctx.FormValue("page"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	tagId := libs.ParseInt(ctx.FormValue("tagId"), 0)
	typeId := libs.ParseInt(ctx.FormValue("typeId"), 0)
	orderBy := ctx.FormValue("orderBy")
	title := ctx.FormValue("title")
	relation := ctx.FormValue("relation")

	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "type_id",
				Condition: "=",
				Value:     typeId,
			}, {
				Key:       "status",
				Condition: "=",
				Value:     "published",
			},
		},
		OrderBy:   orderBy,
		Limit:     limit,
		Offset:    offset,
		Relations: models.GetRelations(relation, nil),
	}

	s.Fields = append(s.Fields, models.GetSearche("title", title))

	articles, count, err := models.GetAllArticles(s, tagId)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	transform := articlesTransform(articles)
	_, _ = ctx.JSON(ApiResource(200, map[string]interface{}{"items": transform, "total": count, "limit": limit}, "操作成功"))
}

/**
* @api {get} /admin/articles 获取所有的文章
* @apiName 获取所有的文章
* @apiGroup Articles
* @apiVersion 1.0.0
* @apiDescription 获取所有的文章
* @apiSampleRequest /admin/articles
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllArticles(ctx iris.Context) {
	offset := libs.ParseInt(ctx.FormValue("page"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	tagId := libs.ParseInt(ctx.FormValue("tagId"), 0)
	typeId := libs.ParseInt(ctx.FormValue("typeId"), 0)
	orderBy := ctx.FormValue("orderBy")

	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "type_id",
				Condition: "=",
				Value:     typeId,
			},
		},
		OrderBy: orderBy,
		Limit:   limit,
		Offset:  offset,
	}

	articles, count, err := models.GetAllArticles(s, tagId)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	transform := articlesTransform(articles)
	_, _ = ctx.JSON(ApiResource(200, map[string]interface{}{"items": transform, "total": count, "limit": limit}, "操作成功"))
}

func articlesTransform(articles []*models.Article) []*transformer.Article {
	var rs []*transformer.Article
	for _, article := range articles {
		r := articleTransform(article)
		rs = append(rs, r)
	}
	return rs
}

func articleTransform(article *models.Article) *transformer.Article {

	r := &transformer.Article{}
	g := gf.NewTransform(r, article, time.RFC3339)
	_ = g.Transformer()
	var tagNames []string
	if len(article.Tags) > 0 {
		for _, tag := range article.Tags {
			tagNames = append(tagNames, tag.Name)
		}
	}
	r.TagNames = tagNames

	if article.Type != nil {
		transform := ttTransform(article.Type)
		r.Type = *transform
	}

	return r
}
