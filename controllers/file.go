package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"path/filepath"
)

/**
* @api {post} /admin/upload_file 文件上传
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
func UploadFile(ctx iris.Context) {
	f, fh, err := ctx.FormFile("uploadfile")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error while uploading: %s", err.Error())))
		return
	}
	defer f.Close()

	path := filepath.Join("./uploads", fh.Filename)
	_, err = ctx.SaveFormFile(fh, path)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error while uploading: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, path, "操作成功"))
}
