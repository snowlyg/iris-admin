package middleware

import (
	"IrisYouQiKangApi/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"net/http"
	"time"
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
	if token.Revoked == 1 || token.ExpressIn < time.Now().Unix() {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(models.ApiJson{State: false, Data: "", Msg: "token 已经过期"})
		ctx.Next()
	}

	user := new(models.Users)
	user.ID = token.UserId

	user.GetUserById() //获取 user 信息

	ctx.Values().Set("auth_user_id", user.ID)
	ctx.Values().Set("auth_user_name", user.Name)

	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}
