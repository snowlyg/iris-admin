package models

import (
	"github.com/jinzhu/gorm"
)

type OauthToken struct {
	gorm.Model
	Token     string `gorm:"not null default '' comment('Token') VARCHAR(191)"`
	UserId    uint   `gorm:"not null default '' comment('UserId') VARCHAR(191)"`
	Secret    string `gorm:"not null default '' comment('Secret') VARCHAR(191)"`
	ExpressIn int64  `gorm:"not null default 0 comment('是否是标准库') BIGINT(20)"`
	Revoked   int    `gorm:"not null TINYINT(1)"`
}

type Token struct {
	Token string `json:"access_token"`
}

func init() {
	DB.AutoMigrate(&OauthToken{})
}

/**
 * oauth_token
 * @method OauthTokenCreate
 */
func (ot *OauthToken) OauthTokenCreate() ApiJson {

	DB.Create(ot)
	response := Token{ot.Token}

	return ApiJson{Status: true, Data: response, Msg: "登陆成功"}
}

/**
 * 通过 token 获取 access_token 记录
 * @method GetOauthTokenByToken
 * @param  {[type]}       token string [description]
 */
func GetOauthTokenByToken(token string) OauthToken {
	var ot OauthToken
	DB.Where("token =  ?", token).First(&ot)
	return ot
}

/**
 * 通过 user_id 更新 oauth_token 记录
 * @method UpdateOauthTokenByUserId
 *@param  {[type]}       user  *OauthToken [description]
 */
func UpdateOauthTokenByUserId(user_id uint) (affected int64, err error) {

	DB.Model(&OauthToken{}).Where("revoked = ?", 0).Where("user_id = ?", user_id).Updates(map[string]interface{}{"revoked": 1})

	return
}
