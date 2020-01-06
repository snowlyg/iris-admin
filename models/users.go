package models

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

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

type UserRequest struct {
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
	u := User{}
	if err := Db.Where("username = ?", username).First(&u).Error; err != nil {
		fmt.Printf("UserAdminCheckLoginErr:%s \n ", err)
	}
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

	if err := Db.Preload("Role").First(user).Error; err != nil {
		fmt.Printf("GetUserByIdErr:%s \n ", err)
	}

	return user
}

/**
 * 通过 username 获取 user 记录
 * @method GetUserByUserName
 * @param  {[type]}       user  *User [description]
 */
func GetUserByUserName(username string) *User {
	user := &User{Username: username}

	if err := Db.Preload("Role").First(user).Error; err != nil {
		fmt.Printf("GetUserByUserNameErr:%s \n ", err)
	}

	return user
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func DeleteUserById(id uint) {
	u := new(User)
	u.ID = id

	if err := Db.Delete(u).Error; err != nil {
		fmt.Printf("DeleteUserByIdErr:%s \n ", err)
	}
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
func GetAllUsers(name, orderBy string, offset, limit int) []*User {
	var users []*User

	q := GetAll(name, orderBy, offset, limit).Preload("Role")
	if err := q.Find(&users).Error; err != nil {
		fmt.Printf("GetAllUserErr:%s \n ", err)
		return nil
	}
	return users
}

/**
 * 创建
 * @method CreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateUser(aul *UserRequest) (user *User) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(User)
	user.Username = aul.Username
	user.Password = hash
	user.Name = aul.Name
	user.RoleID = aul.RoleID

	if err := Db.Create(user).Error; err != nil {
		fmt.Printf("CreateUserErr:%s \n ", err)
	}

	return
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateUser(uj *UserRequest, id uint) *User {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(uj.Password, salt)

	user := new(User)
	user.ID = id
	uj.Password = hash

	if err := Db.Model(user).Updates(uj).Error; err != nil {
		fmt.Printf("UpdateUserErr:%s \n ", err)
	}

	return user
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

			oauthToken := new(OauthToken)
			oauthToken.Token = tokenString
			oauthToken.UserId = user.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauthToken.CreatedAt = time.Now()

			response = oauthToken.OauthTokenCreate()
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
func UserAdminLogout(userId uint) bool {
	ot := UpdateOauthTokenByUserId(userId)
	return ot.Revoked
}

/**
*创建系统管理员
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemAdmin(roleId uint) *User {

	aul := new(UserRequest)
	aul.Username = "username"
	aul.Password = "password"
	aul.Name = "name"
	aul.RoleID = roleId

	user := GetUserByUserName(aul.Username)

	if user.ID == 0 {
		fmt.Println("创建账号")
		return CreateUser(aul)
	} else {
		fmt.Println("重复初始化账号")
		return user
	}
}
