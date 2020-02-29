package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/database"
	"github.com/snowlyg/IrisAdminApi/libs"
	"github.com/snowlyg/IrisAdminApi/validates"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null VARCHAR(191)"`
	Username string `gorm:"unique;VARCHAR(191)"`
	Password string `gorm:"not null VARCHAR(191)"`
}

func NewUser(id uint, username string) *User {
	return &User{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: username,
	}
}

func NewUserByStruct(ru *validates.CreateUpdateUserRequest) *User {
	return &User{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: ru.Username,
		Name:     ru.Name,
		Password: libs.HashPassword(ru.Password),
	}
}

func (u *User) GetUserByUsername() {
	IsNotFound(database.GetGdb().Where("username = ?", u.Username).First(u).Error)
}

func (u *User) GetUserById() {
	IsNotFound(database.GetGdb().Where("id = ?", u.ID).First(u).Error)
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func (u *User) DeleteUser() {
	if err := database.GetGdb().Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteUserByIdErr:%s \n ", err))
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
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllUserErr:%s \n ", err))
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
func (u *User) CreateUser(aul *validates.CreateUpdateUserRequest) {
	u.Password = libs.HashPassword(aul.Password)
	if err := database.GetGdb().Create(u).Error; err != nil {
		color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
	}

	addRoles(aul, u)

	return
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *User) UpdateUser(uj *validates.CreateUpdateUserRequest) {
	uj.Password = libs.HashPassword(uj.Password)
	if err := Update(u, uj); err != nil {
		color.Red(fmt.Sprintf("UpdateUserErr:%s \n ", err))
	}

	addRoles(uj, u)
}

func addRoles(uj *validates.CreateUpdateUserRequest, user *User) {
	if len(uj.RoleIds) > 0 {
		userId := strconv.FormatUint(uint64(user.ID), 10)
		if _, err := database.GetEnforcer().DeleteRolesForUser(userId); err != nil {
			color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
		}

		for _, roleId := range uj.RoleIds {
			roleId := strconv.FormatUint(uint64(roleId), 10)
			if _, err := database.GetEnforcer().AddRoleForUser(userId, roleId); err != nil {
				color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
			}
		}
	}
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func (u *User) CheckLogin(password string) (*Token, bool, string) {
	if u.ID == 0 {
		return nil, false, "用户不存在"
	} else {
		if ok := bcrypt.Match(password, u.Password); ok {
			token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				"iat": time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))

			oauthToken := new(OauthToken)
			oauthToken.Token = tokenString
			oauthToken.UserId = u.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauthToken.CreatedAt = time.Now()

			response := oauthToken.OauthTokenCreate()

			return response, true, "登陆成功"
		} else {
			return nil, false, "用户名或密码错误"
		}
	}
}

/**
* 用户退出登陆
* @method UserAdminLogout
* @param  {[type]} ids string [description]
 */
func UserAdminLogout(userId uint) bool {
	ot := OauthToken{}
	ot.UpdateOauthTokenByUserId(userId)
	return ot.Revoked
}
