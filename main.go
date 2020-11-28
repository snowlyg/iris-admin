package main

import (
	"flag"
	"fmt"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"path/filepath"

	"github.com/snowlyg/blog/app"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/seeder"
)

var ConfigPath = flag.String("c", "", "配置路径")
var CasbinConfigPath = flag.String("cmp", "", "casbin 模型规则配置文件路径")
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
		fmt.Fprintf(os.Stderr, "  -cmp <path>\n")
		fmt.Fprintf(os.Stderr, "    设置 casbin 模型规则配置文件路径\n")
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

	// 解析参数
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

	// 初始化 easygorm 配置
	casbinModelPath := filepath.Join(libs.CWD(), "rbac_model.conf")
	if *CasbinConfigPath != "" {
		casbinModelPath = *CasbinConfigPath
	}
	easygorm.Init(&easygorm.Config{
		GormConfig: &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   libs.Config.DB.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
				SingularTable: false,                 // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			},
		},
		Adapter:         libs.Config.DB.Adapter,  // 类型
		Name:            libs.Config.DB.Name,     // 数据库名称
		Username:        libs.Config.DB.User,     // 用户名
		Pwd:             libs.Config.DB.Password, // 密码
		Host:            libs.Config.DB.Host,     // 地址
		Port:            libs.Config.DB.Port,     // 端口
		CasbinModelPath: casbinModelPath,         // casbin 模型规则路径
		Debug:           libs.Config.Debug,
		TablePrefix:     "iris", // casbin 模型表前缀
		Models: []interface{}{
			&models.User{},
			&models.Role{},
			&models.Permission{},
			&models.Article{},
			&models.Config{},
			&models.Tag{},
			&models.Type{},
			&models.Doc{},
			&models.Chapter{},
			&models.ChapterIp{},
			&models.ArticleIp{},
		},
	})

	irisServer := app.NewServer()
	if irisServer == nil {
		panic("Http 初始化失败")
	}

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
		if !irisServer.Status {
			panic(fmt.Sprintf("端口 %d 已被使用", libs.Config.Port))
		}
		irisServer.Stop() // 停止
	}

	err := irisServer.Start()
	if err != nil {
		panic(err)
	}

}
