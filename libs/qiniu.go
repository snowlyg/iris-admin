package libs

import (
	"context"
	"github.com/fatih/color"
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/storage"
)

// Upload 上传七牛云
func Upload(path, key string) (string, string, error) {
	putPolicy := storage.PutPolicy{
		Scope: Config.Qiniu.Bucket + ":" + key,
	}

	mac := auth.New(Config.Qiniu.Accesskey, Config.Qiniu.Secretkey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = true

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}
	//putExtra.NoCrc32Check = true
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, path, &putExtra)
	if err != nil {
		color.Red("formUploader.PutFile error: %+v", err)
		return "", "", err
	}
	return ret.Key, ret.Hash, nil
}
