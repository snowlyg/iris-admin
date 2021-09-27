package perm

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	myzap "github.com/snowlyg/iris-admin/server/zap"
	"go.uber.org/zap"
)

// Request 请求参数
type Request struct {
	BasePermission
}

func (req *Request) Request(ctx iris.Context) error {
	if err := ctx.ReadJSON(req); err != nil {
		myzap.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return g.ErrParamValidate
	}
	return nil
}
