package initdb

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Post("/initdb", Init)
	}
	return module.NewModule("/init", handler)
}
