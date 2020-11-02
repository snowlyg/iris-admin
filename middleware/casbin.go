package middleware

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/snowlyg/IrisAdminApi/libs"
	"net/http"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/models"
)

func New(e *casbin.Enforcer) *Casbin {
	return &Casbin{enforcer: e}
}

func (c *Casbin) ServeHTTP(ctx iris.Context) {
	ctx.StatusCode(http.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)
	token := models.OauthToken{}
	token.GetOauthTokenByToken(value.Raw) //获取 access_token 信息
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		//_, _ = ctx.Writef("token 失效，请重新登录") // 输出到前端
		_, _ = ctx.JSON(libs.ApiResource(401, nil, ""))
		ctx.StopExecution()
		return
	} else {
		check, err := c.Check(ctx.Request(), strconv.FormatUint(uint64(token.UserId), 10))
		if !check {
			_, _ = ctx.JSON(libs.ApiResource(401, nil, err.Error()))
			ctx.StopExecution()
			return
		} else {
			ctx.Values().Set("auth_user_id", token.UserId)
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
