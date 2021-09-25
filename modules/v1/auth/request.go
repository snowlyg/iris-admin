package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"go.uber.org/zap"
)

// LoginRequest 登录请求字段
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (req *LoginRequest) Request(ctx iris.Context) error {
	if err := ctx.ReadJSON(req); err != nil {
		g.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return g.ErrParamValidate
	}
	return nil
}
