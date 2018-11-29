package main

import (
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	"time"
)

type Users struct {
	gorm.Model

	Name             string `gorm:"not null VARCHAR(191)"`
	Username         string `gorm:"VARCHAR(191)"`
	Password         string `gorm:"not null VARCHAR(191)"`
	Confirmed        bool
	IsClient         bool
	IsFrozen         bool
	IsAudit          bool
	IsClientAdmin    bool
	IsWechatVerfiy   bool
	WechatName       string `gorm:"VARCHAR(191)"`
	WechatAvatar     string `gorm:"VARCHAR(191)"`
	Email            string `gorm:"unique VARCHAR(191)"`
	OpenId           string `gorm:"unique VARCHAR(191)"`
	Phone            string `gorm:"unique VARCHAR(191)"`
	Role             Roles
	RoleId           uint
	RememberToken    string
	WechatVerfiyTime time.Time
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func MUserAdminCheckLogin(username string) Users {
	var u Users
	db.Where("username =  ?", username).First(&u)
	return u
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) GetUserById() *Users {
	db.First(user)
	return user
}

/**
 * 通过 id 冻结 user
 * @method FrozenUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) FrozenUserById() bool {
	db.Model(&user).Update("is_frozen", true)
	return user.IsFrozen
}

/**
 * 通过 id 解冻 user
 * @method RefrozenUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) RefrozenUserById() bool {
	db.Model(&user).Update("is_frozen", false)

	return !user.IsFrozen
}

/**
 * 通过 id 设置负责人
 * @method SetAuditUserById
 * @param  {[type]}   user  *Users [description]
 */
func (user *Users) SetAuditUserById() bool {

	db.Model(Users{}).Where("is_audit=?", true).Updates(map[string]interface{}{"is_audit": false})
	db.Model(&user).Update("is_audit", true)

	return user.IsAudit
}

/**
 * 通过 id 设置负责人
 * @method SetAuditUserById
 * @param  {[type]}   user  *Users [description]
 */
func (user *Users) DeleteUserById() {
	db.Delete(&user)
}

/**
 * 获取所有的账号
 * @method GetAllUsers
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MGetAllUsers(kw string, cp int, mp int) (users []*Users) {
	if len(kw) > 0 {
		db.Model(Users{}).Where(" is_client = ?", 0).Where("name=?", kw).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)
	}
	db.Model(Users{}).Where(" is_client = ?", 0).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)

	return
}

/**
 * 获取所有的客户联系人
 * @method GetAllClients
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MGetAllClients(kw string, cp int, mp int) (users []*Users) {

	if len(kw) > 0 {
		db.Model(Users{}).Where(" is_client = ?", 1).Where("name=?", kw).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)
	}
	db.Model(Users{}).Where(" is_client = ?", 1).Offset(cp - 1).Limit(mp).Preload("Role").Find(&users)

	return
}

/**
 * 获取所有的客户联系人
 * @method GetAllClients
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MCreateUser(aul *AdminUserLogin) (user *Users) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = &Users{
		Username: aul.Username,
		Password: string(hash),
		Name:     aul.Name,
		Phone:    aul.Phone,
	}

	db.Create(user)

	return
}

/**
 * 获取所有的用户联系人数量
 * @method GetClientCounts
 * @return  {[type]} count int    [description]
 */
func MGetClientCounts() (count int) {
	db.Model(&Users{}).Where("is_client = ?", 1).Count(&count)
	return
}
