package main

import (
	"IrisApiProject/config"
	"IrisApiProject/database"
	"IrisApiProject/models"
	"IrisApiProject/routes"
	"github.com/betacraft/yaag/yaag"
	"github.com/kataras/iris/v12"
)

/**
 * 初始化 iris app
 * @method NewApp
 * @return  {[type]}  api      *iris.Application  [iris app]
 */
func newApp() (api *iris.Application) {
	api = iris.New()

	//同步模型数据表
	//如果模型表这里没有添加模型，单元测试会报错数据表不存在。
	//因为单元测试结束，会删除数据表
	database.DB.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = database.DB.Close()
	})

	routes.Register(api)

	appName := config.Conf.Get("app.name").(string)
	appDoc := config.Conf.Get("app.doc").(string)
	appUrl := config.Conf.Get("app.url").(string)
	//api 文档配置
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: appName,
		DocPath:  appDoc + "/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": appUrl,
			"Staging":    "",
		},
	})

	//初始化系统 账号 权限 角色
	models.CreateSystemData()

	return
}

func main() {
	app := newApp()
	app.RegisterView(iris.HTML("resources", ".html"))
	app.HandleDir("/static", "resources/static") // 设置静态资源

	addr := config.Conf.Get("app.addr").(string)
	_ = app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed))
}
