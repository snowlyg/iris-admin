package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/betacraft/yaag/yaag"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	gf "github.com/snowlyg/gotransformer"
	"github.com/snowlyg/iris-base-rabc/files"
	"github.com/snowlyg/iris-base-rabc/middleware"
	"github.com/snowlyg/iris-base-rabc/models"
	"github.com/snowlyg/iris-base-rabc/routes"
	"github.com/snowlyg/iris-base-rabc/transformer"
	"github.com/snowlyg/iris-base-rabc/validates"
)

func main() {
	f := newLogFile()
	defer f.Close()

	Sc = iris.TOML("./config/conf.tml") // 加载配置文件
	rc := getSysConf()                  //格式化配置文件 other 数据

	api := NewApp(rc)
	//api.Logger().SetOutput(f) //记录日志
	err := api.Run(iris.Addr(rc.App.Port), iris.WithConfiguration(Sc))
	if err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}
