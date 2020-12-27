package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/application/libs"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/libs/response"
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
	f, fh, err := ctx.FormFile("file")
	defer f.Close()
	if err != nil {
		logging.ErrorLogger.Errorf(fmt.Sprintf("Error while uploading: %s", err.Error()))
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	fns := strings.Split(fh.Filename, ".")
	if len(fns) != 2 {
		logging.ErrorLogger.Errorf("Error while uploading: 请上传正确的文件")
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	filename := fmt.Sprintf("%s_%d.%s", libs.MD5(fns[0]), time.Now().Unix(), fns[1])
	path := filepath.Join(libs.CWD(), "uploads", "images")
	err = libs.EnsureDir(path)
	_, err = ctx.SaveFormFile(fh, filepath.Join(path, filename))
	if err != nil {
		logging.ErrorLogger.Errorf("Error while SaveFormFile ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	imageHost := fmt.Sprintf("http://%s:%d", ctx.Domain(), libs.Config.Port)
	qiniuUrl := ""
	path = filepath.Join("uploads", "images", filename)
	if libs.Config.Qiniu.Enable {
		var key string
		var hash string
		key, hash, err = libs.Upload(filepath.Join(libs.CWD(), path), filename)
		if err != nil {
			logging.ErrorLogger.Errorf("图片上传七牛云失败 ", err)
			ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
			return
		}

		logging.DebugLogger.Debugf("key:%s ", key, " hash:%s", hash)

		if key != "" {
			qiniuUrl = fmt.Sprintf("%s/%s", libs.Config.Qiniu.Host, key)
		}
	}

	imageUrl := fmt.Sprintf("%s/%s", imageHost, path)
	ctx.JSON(response.NewResponse(response.NoErr.Code, map[string]string{
		"local": imageUrl,
		"qiniu": qiniuUrl,
	}, response.NoErr.Msg))
}
