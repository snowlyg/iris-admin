package initdb

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web/validate"
	"go.uber.org/zap"
)

// InitDB 初始化项目接口
func Init(ctx iris.Context) {
	req := Request{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			g.ZAPLOG.Error("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	err := InitDB(req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// Check 检测是否需要初始化项目
func Check(ctx iris.Context) {
	if database.Instance() == nil {
		ctx.JSON(g.Response{Code: g.NeedInitErr.Code, Data: iris.Map{
			"needInit": true,
		}, Msg: str.Join(g.NeedInitErr.Msg, ":数据库初始化失败")})
		return
	} else if g.CONFIG.System.CacheType == "redis" && g.CACHE == nil {
		ctx.JSON(g.Response{Code: g.NeedInitErr.Code, Data: iris.Map{
			"needInit": true,
		}, Msg: str.Join(g.NeedInitErr.Msg, ":缓存驱动初始化失败")})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: iris.Map{
		"needInit": false,
	}, Msg: g.NoErr.Msg})
}
