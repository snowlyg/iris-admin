package static

import (
	"path/filepath"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party 静态资源模块
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.HandleDir("/", iris.Dir(filepath.Join(dir.GetCurrentAbPath(), "static")))
	}
	return module.NewModule("/static", handler)
}
