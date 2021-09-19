package file

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"go.uber.org/zap"
)

var (
	ErrEmptyErr = errors.New("请上传正确的文件")
)

func UploadFile(ctx iris.Context, fh *multipart.FileHeader) (iris.Map, error) {
	filename, err := GetFileName(fh.Filename)
	if err != nil {

		return nil, err
	}
	path := filepath.Join(dir.GetCurrentAbPath(), "static", "upload", "images")
	err = dir.InsureDir(path)
	if err != nil {
		g.ZAPLOG.Error("文件上传失败", zap.String("dir.InsureDir", err.Error()))
		return nil, err
	}
	_, err = ctx.SaveFormFile(fh, filepath.Join(path, filename))
	if err != nil {
		g.ZAPLOG.Error("文件上传失败", zap.String("ctx.SaveFormFile", "保存文件到本地"))
		return nil, err
	}

	qiniuUrl := ""
	path = GetPath(filename)
	// if libs.Config.Qiniu.Enable {
	// 	var key string
	// 	var hash string
	// 	key, hash, err = libs.Upload(filepath.Join(libs.CWD(), path), filename)
	// 	if err != nil {
	// 		g.ZAPLOG.Error("文件上传失败", zap.String("ctx.SaveFormFile", "图片上传七牛云失败"))
	// 		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
	// 		return
	// 	}

	// 	g.ZAPLOG.Debug("文件上传失败", zap.String("key", key), zap.String("hash", hash))

	// 	if key != "" {
	// 		qiniuUrl = fmt.Sprintf("%s/%s", libs.Config.Qiniu.Host, key)
	// 	}
	// }

	return iris.Map{"local": path, "qiniu": qiniuUrl}, nil
}

func GetFileName(name string) (string, error) {
	fns := strings.Split(strings.TrimLeft(name, "./"), ".")
	if len(fns) != 2 {
		g.ZAPLOG.Error("文件上传失败", zap.String("trings.Split", name))
		return "", ErrEmptyErr
	}
	ext := fns[1]
	md5, err := dir.MD5(name)
	if err != nil {
		g.ZAPLOG.Error("文件上传失败", zap.String("dir.MD5", err.Error()))
		return "", err
	}
	return str.Join(md5, ".", ext), nil
}

func GetPath(filename string) string {
	return filepath.Join("upload", "images", filename)
}
