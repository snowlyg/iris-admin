package middleware

import (
	"IrisYouQiKangApi/controllers"
	"IrisYouQiKangApi/models"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"net/http"
	"time"
)

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},

		SigningMethod: jwt.SigningMethodHS256,
	})

}

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
		ctx.JSON(controllers.ApiJson{Status: false, Data: "", Msg: "token 已经过期"})
		ctx.Next()

		return
	}

	user := new(models.Users)
	user.ID = token.UserId

	user.GetUserById() //获取 user 信息

	ctx.Values().Set("auth_user_id", user.ID)
	ctx.Values().Set("auth_user_name", user.Name)

	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}
