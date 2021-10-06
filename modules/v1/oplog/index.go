package oplog

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
	"github.com/snowlyg/iris-admin/server/operation"
)

// Party 操作日志
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.MultiHandler(), operation.OperationRecord(), middleware.Casbin())
	}
	return module.NewModule("/oplog", handler)
}
