package initdb

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	myzap "github.com/snowlyg/iris-admin/server/zap"
	"go.uber.org/zap"
)

type Request struct {
	Sql       Sql    `json:"sql"`
	SqlType   string `json:"sqlType" validate:"required"`
	Cache     Cache  `json:"cache"`
	CacheType string `json:"cacheType"  validate:"required"`
	Level     string `json:"level"` // debug,release,test
	Addr      string `json:"addr"`
}

func (req *Request) Request(ctx iris.Context) error {
	if err := ctx.ReadJSON(req); err != nil {
		myzap.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return g.ErrParamValidate
	}
	return nil
}

type Sql struct {
	Host     string `json:"host"  validate:"required"`
	Port     string `json:"port"  validate:"required"`
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password"  validate:"required"`
	DBName   string `json:"dbName" validate:"required"`
	LogMode  bool   `json:"logMode"`
}
type Cache struct {
	Host     string `json:"host"  validate:"required"`
	Port     string `json:"port"  validate:"required"`
	Password string `json:"password"`
	PoolSize int    `json:"poolSize"`
	DB       int    `json:"db"`
}
