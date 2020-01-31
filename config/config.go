package config

import (
	"sync"
	"time"

	"IrisAdminApi/transformer"
	"github.com/kataras/iris/v12"
	gf "github.com/snowlyg/gotransformer"
)

type config struct {
	Cf  *transformer.Conf
	Isc iris.Configuration
}

var cfg *config
var once sync.Once

func getConfig() *config {
	once.Do(func() {
		isc := iris.TOML("./config/conf.tml") // 加载配置文件
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

		cf := &transformer.Conf{
			App:      app,
			Mysql:    db,
			Mongodb:  mongodb,
			Redis:    redis,
			Sqlite:   sqlite,
			TestData: testData,
		}

		cfg = &config{Cf: cf, Isc: isc}
	})

	return cfg
}

func GetIrisConf() iris.Configuration {
	return getConfig().Isc
}

func GetTfConf() *transformer.Conf {
	return getConfig().Cf
}
