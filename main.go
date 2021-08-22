package main

import (
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper"
	"github.com/snowlyg/iris-admin/server/web"
)

var Version = "master"

func main() {
	viper.Init()
	webServer := web.Init()
	webServer.SetAddr(g.CONFIG.System.Addr)
	webServer.Run()

}
