package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Users struct {
	gorm.Model

	Name             string `gorm:"not null VARCHAR(191)"`
	Username         string `gorm:"VARCHAR(191)"`
	Password         string `gorm:"not null VARCHAR(191)"`
	Confirmed        int    `gorm:"not null default 0 TINYINT(1)"`
	IsClient         int    `gorm:"not null default 0 TINYINT(1)"`
	IsFrozen         int    `gorm:"not null default 0 TINYINT(1)"`
	IsAudit          int    `gorm:"not null default 0 TINYINT(1)"`
	IsClientAdmin    int    `gorm:"not null default 0 TINYINT(1)"`
	WechatName       string `gorm:"VARCHAR(191)"`
	WechatAvatar     string `gorm:"VARCHAR(191)"`
	Email            string `gorm:"unique VARCHAR(191)"`
	OpenId           string `gorm:"unique VARCHAR(191)"`
	WechatVerfiyTime time.Time
	IsWechatVerfiy   int    `gorm:"not null default 0 TINYINT(1)"`
	Phone            string `gorm:"unique VARCHAR(191)"`
	Role             Roles
	RoleId           uint
	RememberToken    string `gorm:"VARCHAR(100)"`
}

func init() {
	DB.AutoMigrate(&Users{})
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func UserAdminCheckLogin(username string) Users {
	var u Users
	DB.Where("username =  ?", username).First(&u)
	return u
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) GetUserById() (has bool, err error) {
	DB.First(user)
	return
}

/**
 * 获取所有的账号
 * @method GetAllUsers
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func GetAllUsers(kw string, cp int, mp int) (aj ApiJson) {
	users := make([]Users, 0)
	if len(kw) > 0 {
		DB.Model(Users{}).Where(" is_client = ?", 0).Where("name=?", kw).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)
	}
	DB.Model(Users{}).Where(" is_client = ?", 0).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)

	auts := TransFormUsers(users)

	aj.State = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

/**
 * 获取所有的客户联系人
 * @method GetAllClients
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func GetAllClients(kw string, cp int, mp int) (aj ApiJson) {
	users := make([]Users, 0)
	if len(kw) > 0 {
		DB.Model(Users{}).Where(" is_client = ?", 1).Where("name=?", kw).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)
	}
	DB.Model(Users{}).Where(" is_client = ?", 1).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)

	auts := TransFormUsers(users)

	aj.State = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}
