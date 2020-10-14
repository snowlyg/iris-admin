package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/libs"
	"path/filepath"
	"strings"
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
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error while uploading: %s", err.Error())))
		return
	}
	defer f.Close()

	fns := strings.Split(fh.Filename, ".")
	if len(fns) != 2 {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, "Error while uploading: 请上传正确的文件"))
		return
	}

	filename := fmt.Sprintf("%s.%s", libs.MD5(fns[0]), fns[1])
	path := filepath.Join("./uploads/images", filename)
	_, err = ctx.SaveFormFile(fh, path)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error while SaveFormFile: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	imageUrl := fmt.Sprintf("http://%s:%d/%s", config.Config.Host, config.Config.Port, path)
	_, _ = ctx.JSON(ApiResource(200, imageUrl, "操作成功"))
}
