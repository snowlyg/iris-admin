package main

import (
	"github.com/snowlyg/iris-admin/modules/debug"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/viper"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/zap"
)

var Version = "2.0"

func main() {
	viper.Init()
	zap.Init()
	cache.Init()
	webServer := web.Init()
	webServer.AddModule(debug.Party())
	webServer.Run()
}
