package models

import (
	"IrisYouQiKangApi/system"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	Frozen   = 1
	ReFrozen = 0
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
	system.DB.AutoMigrate(&Users{})
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func UserAdminCheckLogin(username string) Users {
	var u Users
	system.DB.Where("username =  ?", username).First(&u)
	return u
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) GetUserById() (aj ApiJson) {

	system.DB.First(user)
	us := []Users{*user}
	tu := TransFormUsers(us)[0]

	aj.Status = true
	aj.Data = tu
	aj.Msg = "操作成功"

	return
}

/**
 * 通过 id 冻结 user
 * @method FrozenUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) FrozenUserById() (aj ApiJson) {
	system.DB.Model(&user).Update("is_frozen", Frozen)

	if user.IsFrozen == Frozen {
		aj.Status = true
		aj.Msg = "操作成功"
	} else {
		aj.Status = false
		aj.Msg = "操作失败"
	}

	return
}

/**
 * 通过 id 解冻 user
 * @method RefrozenUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) RefrozenUserById() (aj ApiJson) {
	system.DB.Model(&user).Update("is_frozen", ReFrozen)

	if user.IsFrozen == ReFrozen {
		aj.Status = true
		aj.Msg = "操作成功"
	} else {
		aj.Status = false
		aj.Msg = "操作失败"
	}

	return
}

/**
 * 通过 id 设置负责人
 * @method SetAuditUserById
 * @param  {[type]}   user  *Users [description]
 */
func (user *Users) SetAuditUserById() (aj ApiJson) {

	system.DB.Model(Users{}).Where("is_audit=?", 1).Updates(map[string]interface{}{"is_audit": 0})
	system.DB.Model(&user).Update("is_audit", 1)

	if user.IsAudit == 1 {
		aj.Status = true
		aj.Msg = "操作成功"
	} else {
		aj.Status = false
		aj.Msg = "操作失败"
	}

	return
}

/**
 * 通过 id 设置负责人
 * @method SetAuditUserById
 * @param  {[type]}   user  *Users [description]
 */
func (user *Users) DeleteUserById() (aj ApiJson) {
	system.DB.Delete(&user)

	aj.Status = true
	aj.Msg = "操作成功"

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
		system.DB.Model(Users{}).Where(" is_client = ?", 0).Where("name=?", kw).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)
	}
	system.DB.Model(Users{}).Where(" is_client = ?", 0).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)

	auts := TransFormUsers(users)

	aj.Status = true
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
		system.DB.Model(Users{}).Where(" is_client = ?", 1).Where("name=?", kw).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)
	}
	system.DB.Model(Users{}).Where(" is_client = ?", 1).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)

	auts := TransFormUsers(users)

	aj.Status = true
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
func CreateUser(aul *AdminUserLogin) (aj ApiJson) {
	salt, e := bcrypt.Salt(10)
	hash, e := bcrypt.Hash(aul.Password, salt)
	if e != nil {

	}

	user := Users{
		Username: aul.Username,
		Password: string(hash),
		Name:     aul.Name,
		Phone:    aul.Phone,
	}

	system.DB.Create(&user)

	us := []Users{user}
	tu := TransFormUsers(us)[0]

	aj.Status = true
	aj.Data = tu
	aj.Msg = "操作成功"

	return
}
