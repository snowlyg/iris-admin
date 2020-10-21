package controllers

import (
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
	article, err := models.GetPublishedArticleById(id)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
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
	article, err := models.GetArticleById(id)
	if err != nil {
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
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

	err := validates.Validate.Struct(*article)
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

	err = article.CreateArticle()
	if err != nil {
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	if article.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, nil, "操作失败"))
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

	err := validates.Validate.Struct(*article)
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
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
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
	offset := libs.ParseInt(ctx.FormValue("offset"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	searchStr := ctx.FormValue("searchStr")
	orderBy := ctx.FormValue("orderBy")

	articles, count, err := models.GetAllArticles(searchStr, orderBy, "published", offset, limit)
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
	offset := libs.ParseInt(ctx.FormValue("offset"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	searchStr := ctx.FormValue("searchStr")
	orderBy := ctx.FormValue("orderBy")

	articles, count, err := models.GetAllArticles(searchStr, orderBy, "", offset, limit)
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
