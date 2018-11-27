package logic

import (
	"IrisYouQiKangApi/models"
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
func UserAdminCheckLogin(username, password string) (response models.Token, status bool, msg string) {
	user := models.UserAdminCheckLogin(username)
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

			oauth_token := new(models.OauthToken)
			oauth_token.Token = tokenString
			oauth_token.UserId = user.ID
			oauth_token.Secret = "secret"
			oauth_token.Revoked = 0
			oauth_token.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauth_token.CreatedAt = time.Now()

			return oauth_token.OauthTokenCreate()

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
func UserAdminLogout(user_id uint) (json models.ApiJson) {

	models.UpdateOauthTokenByUserId(user_id)

	json.Data = true
	json.Msg = "操作成功"

	return
}

///**
// * 删除用户
// * @method UserAdminDele
// * @param  {[type]} ids string [description]
// */
//func UserAdminDele(ids string) models.ApiJson {
//	idsArr := strings.Split(ids, ",")
//	length := len(idsArr)
//	if length > 0 {
//		var idsInt = make([]int, length, length)
//		for i, id := range idsArr {
//			idsInt[i] =system.Tools.ParseInt(id, 0)
//		}
//		return models.UserAdminDele(idsInt)
//	} else {
//		return models.ApiJson{Status: false, Msg: "id is error"}
//	}
//}
