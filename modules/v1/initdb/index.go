package initdb

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party 初始化模块
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Post("/initdb", Init)
		index.Get("/checkdb", Check)
	}
	return module.NewModule("/init", handler)
}
