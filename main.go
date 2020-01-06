package main

import (
	"fmt"
	"time"

	"IrisApiProject/models"
	"IrisApiProject/routes"
	"IrisApiProject/transformer"
	"github.com/betacraft/yaag/yaag"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	gf "github.com/snowlyg/gotransformer"
)

var Sc iris.Configuration

func main() {
	api := iris.New()
	api.Logger().SetLevel("error")
	api.Use(recover.New())
	api.Use(logger.New())
	api.RegisterView(iris.HTML("resources", ".html"))
	api.HandleDir("/static", "resources/static") // 设置静态资源

	Sc = iris.TOML("./config/conf.tml")
	rc := getSysConf()

	models.Register(rc)
	//同步模型数据表
	//如果模型表这里没有添加模型，单元测试会报错数据表不存在。
	//因为单元测试结束，会删除数据表
	models.Db.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)
	iris.RegisterOnInterrupt(func() {
		_ = models.Db.Close()
	})

	//api 文档配置
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: "IrisAdminApi",
		DocPath:  "./apidoc/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": "http://localhost",
			"Staging":    "",
		},
	})

	routes.Register(api)      //注册路由
	models.CreateSystemData() //初始化系统 账号 权限 角色

	err := api.Run(iris.Addr(rc.App.Port), iris.WithConfiguration(Sc))
	if err != nil {
		fmt.Println(err)
	}
}

// 获取配置信息
func getSysConf() *transformer.Conf {

	app := transformer.App{}
	g := gf.NewTransform(&app, Sc.Other["App"], time.RFC3339)
	_ = g.Transformer()

	db := transformer.Database{}
	g.OutputObj = &db
	g.InsertObj = Sc.Other["Database"]
	_ = g.Transformer()

	mongodb := transformer.Mongodb{}
	g.OutputObj = &mongodb
	g.InsertObj = Sc.Other["Mongodb"]
	_ = g.Transformer()

	redis := transformer.Redis{}
	g.OutputObj = &redis
	g.InsertObj = Sc.Other["Redis"]
	_ = g.Transformer()

	sqlite := transformer.Sqlite{}
	g.OutputObj = &sqlite
	g.InsertObj = Sc.Other["Sqlite"]
	_ = g.Transformer()

	testData := transformer.TestData{}
	g.OutputObj = &testData
	g.InsertObj = Sc.Other["TestData"]
	_ = g.Transformer()

	cf := &transformer.Conf{}
	cf.App = app
	cf.Database = db
	cf.Mongodb = mongodb
	cf.Redis = redis
	cf.Sqlite = sqlite
	cf.TestData = testData

	return cf
}
