package main

import (
	"fmt"
	"os"
	"time"

	"IrisAdminApi/config"
	"IrisAdminApi/database"
	"IrisAdminApi/files"
	"IrisAdminApi/models"
	"IrisAdminApi/routes"
	"github.com/betacraft/yaag/yaag"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
)

func newLogFile() *os.File {
	path := "./logs/"
	_ = files.CreateFile(path)
	filename := path + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		color.Red(fmt.Sprintf("日志记录出错: %v", err))
	}

	return f
}

func NewApp() *iris.Application {
	api := iris.New()
	api.Logger().SetLevel(config.GetAppLoggerLevel())

	api.RegisterView(iris.HTML("resources", ".html"))
	//api.HandleDir("/admin", "resources/admin")
	//api.HandleDir("/static", "resources/app/static")

	db := database.GetGdb()
	db.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = db.Close()
	})

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware. //api 文档配置
		On:       true,
		DocTitle: config.GetAppName(),
		DocPath:  "./resources/apiDoc/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": config.GetAppUrl(),
			"Staging":    "",
		},
	})

	routes.App(api) //注册 app 路由
	//routes.Admin(api) //注册 admin 路由

	return api
}

func main() {
	f := newLogFile()
	defer f.Close()

	api := NewApp()
	api.Logger().SetOutput(f) //记录日志
	err := api.Run(iris.Addr(config.GetAppUrl()), iris.WithConfiguration(config.GetIrisConf()))
	if err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}
