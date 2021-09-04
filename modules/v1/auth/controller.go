package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/multi"
)

// Login 登录
func Login(ctx iris.Context) {
	token, err := GetAccessToken(multi.GetUserId(ctx))
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: iris.Map{"accessToken": token}, Msg: g.NoErr.Msg})
}
