package file

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/operation"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
)

// Party 上传文件模块
func Party() func(index iris.Party) {
	return func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.MultiHandler(), operation.OperationRecord(), casbin.Casbin())
		index.Post("/", iris.LimitRequestBodySize(web_iris.CONFIG.MaxSize+1<<20), Upload).Name = "上传文件"
	}
}
