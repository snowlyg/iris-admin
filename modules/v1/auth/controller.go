package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
)

// Login 登录
// - LoginRequest 登录需要提交 uesrname 和 password 字段到接口
// - validate.ValidRequest 验证接口提交参数，需要在 LoginRequest 的字段设置 validate:"required"
// - GetAccessToken 生成验证 token
func Login(ctx iris.Context) {
	req := &LoginRequest{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	token, err := GetAccessToken(req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: iris.Map{"accessToken": token}, Msg: g.NoErr.Msg})
}
