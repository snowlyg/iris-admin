package models

import (
	"IrisApiProject/config"
	"IrisApiProject/database"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null VARCHAR(191)"`
	Username string `gorm:"unique;VARCHAR(191)"`
	Password string `gorm:"not null VARCHAR(191)"`
	RoleID   uint
	Role     Role
}

type UserJson struct {
	Username string `json:"username" validate:"required,gte=2,lte=50"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required,gte=2,lte=50"`
	RoleID   uint   `json:"role_id" validate:"required"`
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func UserAdminCheckLogin(username string) User {
	var u User
	database.DB.Where("username =  ?", username).First(&u)
	return u
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *User [description]
 */
func GetUserById(id uint) *User {
	user := new(User)
	user.ID = id

	database.DB.Preload("Role").First(user)

	return user
}

/**
 * 通过 username 获取 user 记录
 * @method GetUserByUserName
 * @param  {[type]}       user  *User [description]
 */
func GetUserByUserName(username string) *User {
	user := &User{Username: username}

	database.DB.Preload("Role").First(user)

	return user
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func DeleteUserById(id uint) {
	u := new(User)
	u.ID = id

	database.DB.Delete(u)
}

/**
 * 获取所有的账号
 * @method GetAllUser
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllUsers(name, username, orderBy string, offset, limit int) (users []*User) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name
	searchKeys["username"] = username
	database.GetAll(searchKeys, orderBy, offset, limit).Preload("Role").Find(&users)
	return
}

/**
 * 创建
 * @method CreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateUser(aul *UserJson) (user *User) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(User)
	user.Username = aul.Username
	user.Password = string(hash)
	user.Name = aul.Name
	user.RoleID = aul.RoleID

	database.DB.Create(user)

	return
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateUser(aul *UserJson, id uint) *User {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	var user User
	database.DB.First(&user, id)

	data := map[string]interface{}{"username": aul.Username, "password": string(hash), "name": aul.Name, "role_id": aul.RoleID}

	fmt.Println(aul)

	database.DB.Model(&user).Updates(data)

	return &user
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func CheckLogin(username, password string) (response Token, status bool, msg string) {
	user := UserAdminCheckLogin(username)
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
func UserAdminLogout(user_id uint) bool {
	ot := UpdateOauthTokenByUserId(user_id)
	return ot.Revoked
}

/**
*创建系统管理员
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemAdmin(role_id uint) *User {

	aul := new(UserJson)
	aul.Username = config.Conf.Get("test.LoginUserName").(string)
	aul.Password = config.Conf.Get("test.LoginPwd").(string)
	aul.Name = config.Conf.Get("test.LoginName").(string)
	aul.RoleID = role_id

	user := GetUserByUserName(aul.Username)

	if user.ID == 0 {
		fmt.Println("创建账号")
		return CreateUser(aul)
	} else {
		fmt.Println("重复初始化账号")
		return user
	}
}
