//go:generate go-bindata -prefix "./www/dist" -fs  ./www/dist/...
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/snowlyg/IrisAdminApi/libs"
	"github.com/snowlyg/IrisAdminApi/seeder"
	"github.com/snowlyg/IrisAdminApi/web_server"
)

var ConfigPath = flag.String("c", "", "配置路径")
var PrintVersion = flag.Bool("v", false, "打印版本号")
var SeederData = flag.Bool("s", false, "填充基础数据")
var SyncPerms = flag.Bool("p", true, "同步权限")
var PrintRouter = flag.Bool("r", false, "打印路由列表")
var Version = "master"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  -c <path>\n")
		fmt.Fprintf(os.Stderr, "    设置配置文件路径\n")
		fmt.Fprintf(os.Stderr, "  -v <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    打印版本号\n")
		fmt.Fprintf(os.Stderr, "  -s <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    填充基础数据\n")
		fmt.Fprintf(os.Stderr, "  -p <true or false> 默认为: true\n")
		fmt.Fprintf(os.Stderr, "    同步权限\n")
		fmt.Fprintf(os.Stderr, "  -r <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    打印路由列表\n")
		fmt.Fprintf(os.Stderr, "\n")
		//flag.PrintDefaults()
	}

	flag.Parse()

	log.Println(fmt.Sprintf(` 
============================================
   ___  ___ ___ ___ ___ ___   _   ___ ___ 
  / __|/ _ \_ _| _ \_ _/ __| /_\ | _ \_ _|
 | (_ | (_) | ||   /| |\__ \/ _ \|  _/| | 
  \___|\___/___|_|_\___|___/_/ \_\_| |___|

============================================

version: %s`, Version))

	libs.InitConfig(*ConfigPath)

	//irisServer := web_server.NewServer(AssetFile()) 如果需要前端文件
	irisServer := web_server.NewServer(nil)
	if irisServer == nil {
		panic("Http 初始化失败")
	}
	irisServer.NewApp()

	if *PrintVersion {
		fmt.Println(fmt.Sprintf("版本号：%s", Version))
	}

	if *SeederData {
		fmt.Println("填充数据：")
		fmt.Println()
		seeder.Run()
	}

	if *SyncPerms {
		fmt.Println("同步权限：")
		fmt.Println()
		seeder.AddPerm()
	}

	if *PrintRouter {
		fmt.Println("系统权限：")
		fmt.Println()
		routes := irisServer.GetRoutes()
		for _, route := range routes {
			fmt.Println("+++++++++++++++")
			fmt.Println(fmt.Sprintf("名称 ：%s ", route.DisplayName))
			fmt.Println(fmt.Sprintf("路由地址 ：%s ", route.Name))
			fmt.Println(fmt.Sprintf("请求方式 ：%s", route.Act))
			fmt.Println()
		}
	}

	if libs.IsPortInUse(libs.Config.Port) {
		panic(fmt.Sprintf("端口 %d 已被使用", libs.Config.Port))
	}

	err := irisServer.Serve()
	if err != nil {
		panic(err)
	}

}
