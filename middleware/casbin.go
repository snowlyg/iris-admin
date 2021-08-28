package middleware

// import (
// 	"errors"
// 	"fmt"
// 	"net/http"

// 	"github.com/casbin/casbin/v2"
// 	"github.com/iris-contrib/middleware/jwt"
// 	"github.com/kataras/iris/v12"
// 	"github.com/snowlyg/iris-admin/g"
// 	"github.com/snowlyg/iris-admin/modules/user"
// 	casbinServer "github.com/snowlyg/iris-admin/server/casbin"
// 	"go.uber.org/zap"
// )

// func New() *Casbin {
// 	return &Casbin{enforcer: casbinServer.Instance()}
// }

// func (c *Casbin) ServeHTTP(ctx iris.Context) {
// 	jwt, ok := ctx.Values().Get("jwt").(*jwt.Token)
// 	if !ok {
// 		_, _ = ctx.JSON(g.Response{Code: g.AuthErr.Code, Data: nil, Msg: g.AuthErr.Msg})
// 		ctx.StopExecution()
// 		return
// 	}
// 	sess, err := user.Check(jwt.Raw)
// 	if err != nil {
// 		_, _ = ctx.JSON(g.Response{Code: g.AuthErr.Code, Data: nil, Msg: g.AuthErr.Msg})
// 		ctx.StopExecution()
// 		return
// 	}

// 	if sess == nil {
// 		_, _ = ctx.JSON(g.Response{Code: g.AuthExpireErr.Code, Data: nil, Msg: g.AuthExpireErr.Msg})
// 		ctx.StopExecution()
// 		return
// 	} else {
// 		if check, _ := c.Check(ctx.Request(), sess.UserId); !check {
// 			_, _ = ctx.JSON(g.Response{Code: g.AuthActionErr.Code, Data: nil, Msg: fmt.Sprintf("你未拥有当前操作权限，请联系管理员")})
// 			ctx.StopExecution()
// 			return
// 		}
// 	}

// 	ctx.Next()
// }

// // Casbin is the auth services which contains the casbin enforcer.
// type Casbin struct {
// 	enforcer *casbin.Enforcer
// }

// // Check checks the username, request's method and path and
// // returns true if permission grandted otherwise false.
// func (c *Casbin) Check(r *http.Request, userId string) (bool, error) {
// 	method := r.Method
// 	path := r.URL.Path
// 	ok, err := c.enforcer.Enforce(userId, path, method)
// 	if err != nil {
// 		g.ZAPLOG.Error(fmt.Sprintf("验证权限报错：%s-%s-%s", userId, path, method), zap.String("错误", err.Error()))
// 		return false, err
// 	}

// 	g.ZAPLOG.Debug(fmt.Sprintf("权限：%s-%s-%s", userId, path, method))

// 	if !ok {
// 		return ok, errors.New(fmt.Sprintf("你未拥有当前操作权限，请联系管理员"))
// 	}
// 	return ok, nil
// }
