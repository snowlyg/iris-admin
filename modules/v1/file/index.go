package file

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
	"github.com/snowlyg/iris-admin/server/web"
)

// Party 上传文件模块
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Post("/", iris.LimitRequestBodySize(web.CONFIG.MaxSize+1<<20), Upload).Name = "上传文件"
	}
	return module.NewModule("/upload", handler)
}
