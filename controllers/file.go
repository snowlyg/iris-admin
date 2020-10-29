package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"path/filepath"
	"strings"
	"time"
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

	filename := fmt.Sprintf("%s_%d.%s", libs.MD5(fns[0]), time.Now().Unix(), fns[1])
	path := filepath.Join(libs.CWD(), "uploads", "images")
	err = libs.EnsureDir(path)
	_, err = ctx.SaveFormFile(fh, filepath.Join(path, filename))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error while SaveFormFile: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)

	imageHost := fmt.Sprintf("http://%s:%d", ctx.Domain(), libs.Config.Port)
	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "name",
				Condition: "=",
				Value:     "imageHost",
			},
		},
	}
	configByKey, err := models.GetConfig(s)
	if err == nil {
		imageHost = configByKey.Value
	} else {
		fmt.Println(err.Error())
	}

	imageUrl := fmt.Sprintf("%s/%s", imageHost, filepath.Join("uploads", "images", filename))
	_, _ = ctx.JSON(ApiResource(200, imageUrl, "操作成功"))
}
