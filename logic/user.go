package logic

import (
	"IrisYouQiKangApi/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

/**
 * 判断用户是否登录
 * @method UserAdminLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func UserAdminCheckLogin(username, password string) models.ApiJson {
	user := models.UserAdminCheckLogin(username)
	if user.ID == 0 {
		return models.ApiJson{State: false, Msg: "用户不存在"}
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err == nil {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			tokenString, err := token.SignedString([]byte("secret"))

			if err != nil {
				return models.ApiJson{State: true, Msg: err.Error()}
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
			return models.ApiJson{State: false, Msg: "用户名或密码错误"}
		}
	}
}

/**
 * 通过 token 获取 access_token 记录
 * @method GetOauthTokenByToken
 * @param  {[type]}       token string [description]
 */
func GetUserById(id uint) models.ApiJson {

	user := new(models.Users)
	user.ID = id

	if user.ID == 0 {
		return models.ApiJson{State: false, Msg: "用户不存在"}
	}

	has, err := user.GetUserById()

	if err != nil {
		return models.ApiJson{State: false, Data: "", Msg: err.Error()}
	}

	if !has {
		return models.ApiJson{State: false, Data: "", Msg: "没有数据"}
	}

	return models.ApiJson{State: true, Data: user, Msg: "操作成功"}

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
//			idsInt[i] = Tools.ParseInt(id, 0)
//		}
//		return models.UserAdminDele(idsInt)
//	} else {
//		return models.ApiJson{State: false, Msg: "id is error"}
//	}
//}
