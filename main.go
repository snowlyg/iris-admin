package main

import (
	"github.com/snowlyg/iris-admin/modules/debug"
	"github.com/snowlyg/iris-admin/server/viper"
	"github.com/snowlyg/iris-admin/server/web"
)

var Version = "master"

func main() {
	viper.Init()
	webServer := web.Init()
	webServer.AddModule(debug.Party())
	webServer.Run()
}
