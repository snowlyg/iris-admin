package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"IrisAdminApi/files"
	"IrisAdminApi/middleware"
	"IrisAdminApi/models"
	"IrisAdminApi/routes"
	"IrisAdminApi/transformer"
	"github.com/betacraft/yaag/yaag"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	gf "github.com/snowlyg/gotransformer"
)

var Sc iris.Configuration

// 获取路由信息
func getRoutes(api *iris.Application) []*models.PermissionRequest {
	rs := api.APIBuilder.GetRoutes()
	var rrs []*models.PermissionRequest
	for _, s := range rs {
		if !isPermRoute(s) {
			path := strings.Replace(s.Path, ":id", "*", 1)
			rr := &models.PermissionRequest{Name: path, DisplayName: s.Name, Description: s.Name, Act: s.Method}
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

// 过滤非必要权限
func isPermRoute(s *router.Route) bool {
	exceptRouteName := []string{"OPTIONS", "GET", "POST", "HEAD", "PUT", "PATCH"}
	for _, er := range exceptRouteName {
		if strings.Contains(s.Name, er) {
			return true
		}
	}
	return false
}

// 获取配置信息
func getSysConf() *transformer.Conf {
	app := transformer.App{}
	g := gf.NewTransform(&app, Sc.Other["App"], time.RFC3339)
	_ = g.Transformer()

	db := transformer.Mysql{}
	g.OutputObj = &db
	g.InsertObj = Sc.Other["Mysql"]
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

	cf := &transformer.Conf{
		App:      app,
		Mysql:    db,
		Mongodb:  mongodb,
		Redis:    redis,
		Sqlite:   sqlite,
		TestData: testData,
	}

	return cf
}

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

func NewApp(rc *transformer.Conf) *iris.Application {

	api := iris.New()
	api.Logger().SetLevel(rc.App.LoggerLevel)

	api.RegisterView(iris.HTML("resources", ".html"))
	api.HandleDir("/static", "resources/static")

	models.Register(rc) // 数据初始化
	models.Db.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)
	iris.RegisterOnInterrupt(func() {
		_ = models.Db.Close()
	})

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware. //api 文档配置
		On:       true,
		DocTitle: rc.App.Name,
		DocPath:  "./resources/apiDoc/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": rc.App.URl + rc.App.Port,
			"Staging":    "",
		},
	})

	routes.Register(api)                   //注册路由
	middleware.Register(api)               // 中间件注册
	apiRoutes := getRoutes(api)            // 获取路由数据
	models.CreateSystemData(rc, apiRoutes) // 初始化系统数据 管理员账号，角色，权限

	return api
}

func main() {
	f := newLogFile()
	defer f.Close()

	Sc = iris.TOML("./config/conf.tml") // 加载配置文件
	rc := getSysConf()                  //格式化配置文件 other 数据
	api := NewApp(rc)
	api.Logger().SetOutput(f) //记录日志
	err := api.Run(iris.Addr(rc.App.Port), iris.WithConfiguration(Sc))
	if err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}
