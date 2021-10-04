package perm

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
)

// Request 请求参数
type Request struct {
	BasePermission
}

func (req *Request) Request(ctx iris.Context) error {
	if err := ctx.ReadJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}
