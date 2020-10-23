package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
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
	value := ctx.Values().Get("jwt").(*jwt.Token)
	//token := models.OauthToken{}
	//token.GetOauthTokenByToken(value.Raw) //获取 access_token 信息
	//if token.Revoked || token.ExpressIn < time.Now().Unix() {
	//	//_, _ = ctx.Writef("token 失效，请重新登录") // 输出到前端
	//	ctx.StatusCode(http.StatusUnauthorized)
	//	ctx.StopExecution()
	//	return
	//} else if !c.Check(ctx.Request(), strconv.FormatUint(uint64(token.UserId), 10)) {
	//	ctx.StatusCode(http.StatusForbidden) // Status Forbidden
	//	ctx.StopExecution()
	//	return
	//} else {
	//	ctx.Values().Set("auth_user_id", token.UserId)
	//}

	conn := libs.GetRedisClusterClient()
	defer conn.Close()

	sess, err := models.GetRedisSessionV2(conn, value.Raw)
	if err != nil {
		fmt.Println(fmt.Sprintf("get %s redis session error : %s\n", value.Raw, err.Error()))
		models.UserTokenExpired(value.Raw)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.StopExecution()
		return
	}
	if sess == nil {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.StopExecution()
		return
	} else {
		ctx.Values().Set("sess", sess)
	}

	ctx.Next()
}

// Casbin is the auth services which contains the casbin enforcer.
type Casbin struct {
	enforcer *casbin.Enforcer
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
func (c *Casbin) Check(r *http.Request, userId string) bool {
	method := r.Method
	path := r.URL.Path
	ok, err := c.enforcer.Enforce(userId, path, method)
	fmt.Println()
	if err != nil {
		fmt.Println(fmt.Sprintf("验证权限报错：%v;%s-%s-%s", err.Error(), userId, path, method))
		return false
	}
	return ok
}
