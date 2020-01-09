package middleware

import (
	"net/http"
	"time"

	"IrisAdminApi/controllers"
	"IrisAdminApi/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
)

/**
 * 判断 token 是否有效
 * 如果有效 就获取信息并且保存到请求里面
 * @method AuthToken
 * @param  {[type]}  ctx       iris.Context    [description]
 */
func AuthToken(ctx iris.Context) {
	u := ctx.Values().Get("jwt").(*jwt.Token)   //获取 token 信息
	token := models.GetOauthTokenByToken(u.Raw) //获取 access_token 信息
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		ctx.StatusCode(http.StatusUnauthorized)
		_, _ = ctx.JSON(controllers.ApiResource(false, nil, "token 已经过期"))
		//ctx.Next()
		//return
	} else {
		ctx.Values().Set("auth_user_id", token.UserId)
	}

	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}
