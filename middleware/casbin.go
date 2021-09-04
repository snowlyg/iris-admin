package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/iris-admin/g"
	casbinServer "github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

func Casbin() iris.Handler {
	return func(ctx *context.Context) {
		check, err := Check(ctx.Request(), strconv.FormatUint(uint64(multi.GetUserId(ctx)), 10))
		if err != nil || !check {
			_, _ = ctx.JSON(g.Response{Code: g.AuthActionErr.Code, Data: nil, Msg: err.Error()})
			ctx.StopExecution()
			return
		}

		ctx.Next()
	}
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
func Check(r *http.Request, userId string) (bool, error) {
	method := r.Method
	path := r.URL.Path
	ok, err := casbinServer.Instance().Enforce(userId, path, method)
	if err != nil {
		g.ZAPLOG.Error(fmt.Sprintf("验证权限报错：%s-%s-%s", userId, path, method), zap.String("错误:", err.Error()))
		return false, err
	}

	g.ZAPLOG.Debug(fmt.Sprintf("权限：%s-%s-%s", userId, path, method))

	if !ok {
		return ok, errors.New("你未拥有当前操作权限，请联系管理员")
	}
	return ok, nil
}
