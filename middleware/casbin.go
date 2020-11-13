package middleware

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"net/http"
)

func New(e *casbin.Enforcer) *Casbin {
	return &Casbin{enforcer: e}
}

func (c *Casbin) ServeHTTP(ctx iris.Context) {
	ctx.StatusCode(http.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)

	conn := libs.GetRedisClusterClient()
	defer conn.Close()

	sess, err := models.GetRedisSessionV2(conn, value.Raw)
	if err != nil {
		models.UserTokenExpired(value.Raw)
		_, _ = ctx.JSON(libs.ApiResource(401, nil, ""))
		ctx.StopExecution()
		return
	}
	if sess == nil {
		_, _ = ctx.JSON(libs.ApiResource(401, nil, ""))
		ctx.StopExecution()
		return
	} else {
		check, err := c.Check(ctx.Request(), sess.UserId)
		if !check {
			_, _ = ctx.JSON(libs.ApiResource(403, nil, err.Error()))
			ctx.StopExecution()
			return
		} else {
			ctx.Values().Set("sess", sess)
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
		color.Red("验证权限报错：%v;%s-%s-%s", err.Error(), userId, path, method)
		return false, err
	}
	if !ok {
		return ok, errors.New(fmt.Sprintf("你未拥有 %s:%s 操作权限", method, path))
	}
	return ok, nil
}
