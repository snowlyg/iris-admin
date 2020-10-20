package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/libs"
	"github.com/snowlyg/IrisAdminApi/sysinit"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null VARCHAR(50)" json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Username string `gorm:"unique;VARCHAR(150)" json:"password" validate:"required"  comment:"密码"`
	Password string `gorm:"not null VARCHAR(50)" json:"name" validate:"required,gte=2,lte=50"  comment:"名称"`
	Intro    string `gorm:"not null VARCHAR(500)" json:"introduction" comment:"简介"`
	Avatar   string `gorm:"not null VARCHAR(500)" json:"avatar"  comment:"头像"`
	RoleIds  []uint `gorm:"-" json:"role_ids"  validate:"required" comment:"角色"`
}

func NewUser() *User {
	return &User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func GetUserByUsername(username string) (*User, error) {
	u := NewUser()
	err := IsNotFound(sysinit.Db.Where("username = ?", username).First(u).Error)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserById(id uint) (*User, error) {
	u := NewUser()
	err := IsNotFound(sysinit.Db.Where("id = ?", id).First(u).Error)
	if err != nil {
		return nil, err
	}
	return u, nil
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func DeleteUser(id uint) error {
	u := NewUser()
	u.ID = id
	if err := sysinit.Db.Delete(u).Error; err != nil {
		return err
	}
	return nil
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
func (u *User) CreateUser() error {
	u.Password = libs.HashPassword(u.Password)
	if err := sysinit.Db.Create(u).Error; err != nil {
		return err
	}

	addRoles(u)

	return nil
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *User) UpdateUser() error {
	u.Password = libs.HashPassword(u.Password)
	if err := Update(&User{}, u); err != nil {
		return err
	}

	addRoles(u)
	return nil
}

func addRoles(user *User) {
	if len(user.RoleIds) > 0 {
		userId := strconv.FormatUint(uint64(user.ID), 10)
		if _, err := sysinit.Enforcer.DeleteRolesForUser(userId); err != nil {
			color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
		}

		for _, roleId := range user.RoleIds {
			roleId := strconv.FormatUint(uint64(roleId), 10)
			if _, err := sysinit.Enforcer.AddRoleForUser(userId, roleId); err != nil {
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
func (u *User) CheckLogin(password string) (*Token, int64, string) {
	if u.ID == 0 {
		return nil, 400, "用户不存在"
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

			return response, 200, "登陆成功"
		} else {
			return nil, 400, "用户名或密码错误"
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
