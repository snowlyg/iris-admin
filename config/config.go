//此包用于获取配置，
//iris 框架本身的配置处理已经比较完善，
//增加这些方法主要是增加配置使用的灵活性
package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/transformer"
	gf "github.com/snowlyg/gotransformer"
)

var (
	Root = os.Getenv("GOPATH") + "/src/github.com/snowlyg/IrisAdminApi"
	Isc  = iris.TOML(filepath.Join(Root, "config", "conf.tml")) // 加载配置文件
)

func newConfig() *transformer.Conf {
	return getTfConf(Isc)
}

func getTfConf(isc iris.Configuration) *transformer.Conf {
	app := transformer.App{}
	g := gf.NewTransform(&app, isc.Other["App"], time.RFC3339)
	_ = g.Transformer()

	db := transformer.Mysql{}
	g.OutputObj = &db
	g.InsertObj = isc.Other["Mysql"]
	_ = g.Transformer()

	mongodb := transformer.Mongodb{}
	g.OutputObj = &mongodb
	g.InsertObj = isc.Other["Mongodb"]
	_ = g.Transformer()

	redis := transformer.Redis{}
	g.OutputObj = &redis
	g.InsertObj = isc.Other["Redis"]
	_ = g.Transformer()

	sqlite := transformer.Sqlite{}
	g.OutputObj = &sqlite
	g.InsertObj = isc.Other["Sqlite"]
	_ = g.Transformer()

	testData := transformer.TestData{}
	g.OutputObj = &testData
	g.InsertObj = isc.Other["TestData"]
	_ = g.Transformer()

	return &transformer.Conf{
		App:      app,
		Mysql:    db,
		Mongodb:  mongodb,
		Redis:    redis,
		Sqlite:   sqlite,
		TestData: testData,
	}
}

func GetAppName() string {
	return newConfig().App.Name
}

func GetAppUrl() string {
	return newConfig().App.Url
}

func GetAppLoggerLevel() string {
	return newConfig().App.LoggerLevel
}

func GetAppDriverType() string {
	return newConfig().App.DriverType
}

func GetAppCreateSysData() bool {
	return newConfig().App.CreateSysData
}

func GetMysqlConnect() string {
	return newConfig().Mysql.Connect
}

func GetMysqlName() string {
	return newConfig().Mysql.Name
}

func GetMysqlTName() string {
	return newConfig().Mysql.TName
}

func GetMongodbConnect() string {
	return newConfig().Mongodb.Connect
}

func GetSqliteConnect() string {
	return filepath.Join(Root, "tmp", newConfig().Sqlite.Connect)
}

func GetSqliteTConnect() string {
	return filepath.Join(Root, "tmp", newConfig().Sqlite.TConnect)
}

func GetTestDataUserName() string {
	return newConfig().TestData.UserName
}

func GetTestDataName() string {
	return newConfig().TestData.Name
}

func GetTestDataPwd() string {
	return newConfig().TestData.Pwd
}
