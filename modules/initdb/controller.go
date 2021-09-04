package initdb

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/validate"
	"go.uber.org/zap"
)

// InitDB 登录
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
