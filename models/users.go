package models

import (
	"IrisApiProject/database"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
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
	RememberToken    string
	WechatVerfiyTime time.Time
}

type UserJson struct {
	Username string `json:"username" validate:"required,gte=4,lte=50"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required,gte=4,lte=50"`
	Phone    string `json:"phone" validate:"required"`
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func MUserAdminCheckLogin(username string) Users {
	var u Users
	database.DB.Where("username =  ?", username).First(&u)
	return u
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) GetUserById() *Users {
	database.DB.First(user)
	return user
}

/**
 * 通过 id 冻结 user
 * @method FrozenUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) FrozenUserById() bool {
	database.DB.Model(&user).Update("is_frozen", true)
	return user.IsFrozen
}

/**
 * 通过 id 解冻 user
 * @method RefrozenUserById
 * @param  {[type]}       user  *Users [description]
 */
func (user *Users) RefrozenUserById() bool {
	database.DB.Model(&user).Update("is_frozen", false)

	return !user.IsFrozen
}

/**
 * 通过 id 设置负责人
 * @method SetAuditUserById
 * @param  {[type]}   user  *Users [description]
 */
func (user *Users) SetAuditUserById() bool {
	database.DB.Model(Users{}).Where("is_audit=?", true).Updates(map[string]interface{}{"is_audit": false})
	database.DB.Model(&user).Update("is_audit", true)

	return user.IsAudit
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func (user *Users) DeleteUserById() {
	database.DB.Delete(&user)
}

/**
 * 获取所有的账号
 * @method MGetAllUsers
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func MGetAllUsers(name, username, orderBy string, offset, limit int) (users []*Users) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name
	searchKeys["username"] = username
	searchKeys["is_client"] = false

	MGetAll(searchKeys, orderBy, "Role", offset, limit).Find(&users)
	return
}

/**
 * 获取所有的客户联系人
 * @method MGetAllClients
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func MGetAllClients(name, username, orderBy string, offset, limit int) (users []*Users) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name
	searchKeys["username"] = username
	searchKeys["is_client"] = false

	MGetAll(searchKeys, orderBy, "Role", offset, limit).Find(&users)
	return
}

/**
 * 创建
 * @method MCreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MCreateUser(aul *UserJson) (user *Users) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(Users)
	user.Username = aul.Username
	user.Password = string(hash)
	user.Name = aul.Name
	user.Phone = aul.Phone

	database.DB.Create(user)

	return
}

/**
 * 获取所有的客户联系人
 * @method MCreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MUpdateUser(aul *UserJson) (user *Users) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(Users)
	user.Username = aul.Username
	user.Password = string(hash)
	user.Name = aul.Name
	user.Phone = aul.Phone

	database.DB.Update(user)

	return
}

/**
 * 获取所有的用户联系人数量
 * @method MGetClientCounts
 * @return  {[type]} count int    [description]
 */
func MGetClientCounts() (count int) {
	database.DB.Model(&Users{}).Where("is_client = ?", 1).Count(&count)
	return
}

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
	ot := UpdateOauthTokenByUserId(user_id)

	return ot.Revoked
}
