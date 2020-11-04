package controllers

import (
	"fmt"
	"github.com/fatih/color"
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

	ctx.StatusCode(iris.StatusOK)
	f, fh, err := ctx.FormFile("uploadfile")
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error while uploading: %s", err.Error())))
		return
	}
	defer f.Close()

	fns := strings.Split(fh.Filename, ".")
	if len(fns) != 2 {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, "Error while uploading: 请上传正确的文件"))
		return
	}

	filename := fmt.Sprintf("%s_%d.%s", libs.MD5(fns[0]), time.Now().Unix(), fns[1])
	path := filepath.Join(libs.CWD(), "uploads", "images")
	err = libs.EnsureDir(path)
	_, err = ctx.SaveFormFile(fh, filepath.Join(path, filename))
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error while SaveFormFile: %s", err.Error())))
		return
	}

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
		color.Red(err.Error())
	}

	qiniuKey := ""
	path = filepath.Join("uploads", "images", filename)
	if libs.Config.Qiniu.Enable {
		key, hash, err := libs.Upload(path, filename)
		if err != nil {
			_, _ = ctx.JSON(libs.ApiResource(200, nil, err.Error()))
			return
		}

		color.Yellow("key:%s,hash:%s", key, hash)
		if key != "" {
			qiniuKey = key
		}
	}

	imageUrl := fmt.Sprintf("%s/%s", imageHost, path)
	qiniuUrl := fmt.Sprintf("%s/%s", libs.Config.Qiniu.Host, qiniuKey)
	_, _ = ctx.JSON(libs.ApiResource(200, map[string]string{
		"local": imageUrl,
		"qiniu": qiniuUrl,
	}, "操作成功"))
}
