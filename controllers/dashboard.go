package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
)

/**
* @api {post} /admin/dashboard 面板数据
* @apiName 文件上传
* @apiGroup UploadFile
* @apiVersion 1.0.0
* @apiDescription 文件上传
* @apiSampleRequest /admin/upload_file
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func Dashboard(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	var (
		articles int64
		docs     int64
		err      error
	)

	articles, err = models.GetArticleCount()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	articleReads, err := models.GetArticleReads()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	docs, err = models.GetDocCount()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	docReads, err := models.GetDocReads()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}

	dashboard := map[string]int64{
		"articleVisitis": articleReads,
		"articles":       articles,
		"docVisitis":     docReads,
		"docs":           docs,
	}
	_, _ = ctx.JSON(libs.ApiResource(200, dashboard, "操作成功"))
}
