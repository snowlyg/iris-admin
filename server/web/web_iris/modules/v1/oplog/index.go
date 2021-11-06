package oplog

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/operation"
	"github.com/snowlyg/iris-admin/server/web/web_iris/middleware"
)

// Party 操作日志
func Party() func(index iris.Party) {
	return func(index iris.Party) {
		index.Use(middleware.MultiHandler(), operation.OperationRecord(), casbin.Casbin())
	}
}
