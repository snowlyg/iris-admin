package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/snowlyg/blog/libs"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null; type:varchar(60)" json:"name" validate:"required,gte=2,lte=50" comment:"用户名"`
	Username string `gorm:"unique;not null;type:varchar(60)" json:"username" validate:"required,gte=2,lte=50"  comment:"名称"`
	Password string `gorm:"type:varchar(100)" json:"password"  comment:"密码"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"introduction" comment:"简介"`
	Avatar   string `gorm:"type:longText" json:"avatar"  comment:"头像"`
	RoleIds  []uint `gorm:"-" json:"role_ids"  validate:"required" comment:"角色"`
}

type Token struct {
	Token string `json:"token"`
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
	user := new(User)
	if err := IsNotFound(libs.Db.Where("username = ?", username).First(user).Error); err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserById(id uint) (*User, error) {
	u := NewUser()
	err := IsNotFound(libs.Db.Where("id = ?", id).First(u).Error)
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
	u, err := GetUserById(id)
	if err != nil {
		return err
	}
	if u.Username == "username" {
		return errors.New("不能删除管理员")
	}

	if err := libs.Db.Delete(u, id).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteUserByIdErr:%s \n ", err))
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
	q, _ := GetAll(&User{}, name, orderBy, offset, limit)
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
	if err := libs.Db.Create(u).Error; err != nil {
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
func UpdateUserById(id uint, nu *User) error {
	nu.Password = libs.HashPassword(nu.Password)
	if err := Update(&User{}, nu, id); err != nil {
		return err
	}

	addRoles(nu)
	return nil
}

func addRoles(user *User) {
	if len(user.RoleIds) > 0 {
		userId := strconv.FormatUint(uint64(user.ID), 10)
		if _, err := libs.Enforcer.DeleteRolesForUser(userId); err != nil {
			color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
		}

		for _, roleId := range user.RoleIds {
			roleId := strconv.FormatUint(uint64(roleId), 10)
			if _, err := libs.Enforcer.AddRoleForUser(userId, roleId); err != nil {
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
		uid := strconv.FormatUint(uint64(u.ID), 10)
		if isUserTokenOver(uid) {
			return nil, 400, "以达到同时登录设备上限"
		}
		if ok := bcrypt.Match(password, u.Password); ok {
			token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				"iat": time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))

			rsv2 := RedisSessionV2{
				UserId:       uid,
				LoginType:    LoginTypeWeb,
				AuthType:     AuthPwd,
				CreationDate: time.Now().Unix(),
				Scope:        getUserScope("admin"),
			}
			conn := libs.GetRedisClusterClient()
			defer conn.Close()

			if err := rsv2.ToCache(conn, tokenString); err != nil {
				return nil, 400, err.Error()
			}

			if err := rsv2.SyncUserTokenCache(conn, tokenString); err != nil {
				return nil, 400, err.Error()
			}

			return &Token{tokenString}, 200, "登陆成功"
		} else {
			return nil, 400, "用户名或密码错误"
		}
	}
}
