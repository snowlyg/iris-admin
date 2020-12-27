package middleware

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/libs/response"
	"github.com/snowlyg/blog/service/auth"
	"net/http"
)

func New() *Casbin {
	return &Casbin{enforcer: easygorm.EasyGorm.Enforcer}
}

func (c *Casbin) ServeHTTP(ctx iris.Context) {
	ctx.StatusCode(http.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)

	authDriver := auth.NewAuthDriver()
	defer authDriver.Close()
	rsv2, err := auth.Check(authDriver, value.Raw)
	if err != nil {
		authDriver.DelUserTokenCache(value.Raw)
		_, _ = ctx.JSON(response.NewResponse(response.AuthErr.Code, nil, response.AuthErr.Msg))
		ctx.StopExecution()
		return
	}

	if rsv2 == nil {
		_, _ = ctx.JSON(response.NewResponse(response.AuthExpireErr.Code, nil, response.AuthExpireErr.Msg))
		ctx.StopExecution()
		return
	} else {
		if check, _ := c.Check(ctx.Request(), rsv2.UserId); !check {
			_, _ = ctx.JSON(response.NewResponse(response.AuthActionErr.Code, nil, fmt.Sprintf("你未拥有当前操作权限，请联系管理员")))
			ctx.StopExecution()
			return
		} else {
			ctx.Values().Set("sess", rsv2)
		}
	}

	ctx.Next()
}

// Casbin is the auth services which contains the casbin enforcer.
type Casbin struct {
	enforcer *casbin.Enforcer
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
func (c *Casbin) Check(r *http.Request, userId string) (bool, error) {
	method := r.Method
	path := r.URL.Path
	ok, err := c.enforcer.Enforce(userId, path, method)
	if err != nil {
		logging.ErrorLogger.Error("验证权限报错：%v;%s-%s-%s", err.Error(), userId, path, method)
		return false, err
	}

	logging.DebugLogger.Debugf("权限：%s-%s-%s", userId, path, method)

	if !ok {
		return ok, errors.New(fmt.Sprintf("你未拥有当前操作权限，请联系管理员"))
	}
	return ok, nil
}
