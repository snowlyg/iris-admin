package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jameskeane/bcrypt"
	"time"
)

/**
 * 判断用户是否登录
 * @method UserAdminLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func LUserAdminCheckLogin(username, password string) (response Token, status bool, msg string) {
	user := MUserAdminCheckLogin(username)
	if user.ID == 0 {
		msg = "用户不存在"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			tokenString, err := token.SignedString([]byte("secret"))

			if err != nil {
				msg = err.Error()
				return
			}

			oauth_token := new(OauthToken)
			oauth_token.Token = tokenString
			oauth_token.UserId = user.ID
			oauth_token.Secret = "secret"
			oauth_token.Revoked = false
			oauth_token.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauth_token.CreatedAt = time.Now()

			response = oauth_token.OauthTokenCreate()
			status = true
			msg = "登陆成功"

			return

		} else {
			msg = "用户名或密码错误"
			return
		}
	}
}

/**
* 用户退出登陆
* @method UserAdminLogout
* @param  {[type]} ids string [description]
 */
func LUserAdminLogout(user_id uint) bool {
	ot := MUpdateOauthTokenByUserId(user_id)

	return ot.Revoked
}
