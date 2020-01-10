package middleware

import (
	"net/http"
	"time"

	"IrisAdminApi/models"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

/**
 * 判断 token 是否有效
 * 如果有效 就获取信息并且保存到请求里面
 * @method AuthToken
 * @param  {[type]}  ctx       iris.Context    [description]
 */
func AuthToken(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)
	token := models.GetOauthTokenByToken(user.Raw) //获取 access_token 信息
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		_, _ = ctx.Writef("token 失效，请重新登录")
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Next()
		//return
	} else {
		ctx.Values().Set("auth_user_id", token.UserId)
	}

	//foobar := user.Claims.(jwt.MapClaims)
	//for key, value := range foobar {
	//	_, _ = ctx.Writef("%s = %s", key, value)
	//}

	ctx.Next()
}
