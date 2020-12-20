package models

import (
	"errors"
	"fmt"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"strconv"
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

type Avatar struct {
	Avatar string `gorm:"type:longText" json:"avatar" validate:"required" comment:"头像"`
}

type Token struct {
	Token string `json:"token"`
}

// GetUser get user
func GetUser(search *easygorm.Search) (*User, error) {
	user := &User{}
	if err := easygorm.First(user, search); err != nil {
		logging.Err.Errorf("get user err: %+v\n", err)
		return user, err
	}
	return user, nil
}

// DeleteUser del user . if user's username is username ,can't del it.
func DeleteUser(id uint) error {
	user := &User{}
	if err := easygorm.FindById(user, id); err != nil {
		return err
	}
	if user.Username == "username" {
		err := errors.New(fmt.Sprintf("不能删除管理员 : %s \n ", user.Username))
		return err
	}
	if err := easygorm.Egm.Db.Delete(user, id).Error; err != nil {
		logging.Err.Errorf("del user err: %+v\n", err)
		return err
	}
	return nil
}

// GetAllUsers get all users
func GetAllUsers(s *easygorm.Search) ([]*User, int64, error) {
	var users []*User
	count, err := easygorm.Paginate(&User{}, &users, s)
	if err != nil {
		logging.Err.Errorf("get all users err: %+v\n", err)
		return nil, count, err
	}
	return users, count, nil
}

// CreateUser create user
func (u *User) CreateUser() error {
	u.Password = libs.HashPassword(u.Password)
	if err := easygorm.Create(u); err != nil {
		logging.Err.Errorf("create user err: %+v\n", err)
		return err
	}
	if err := addRoles(u); err != nil {
		return err
	}
	return nil
}

// UpdateUserById update user by id
func UpdateUserById(id uint, u *User) error {
	if len(u.Password) > 0 {
		u.Password = libs.HashPassword(u.Password)
	}
	if err := easygorm.Update(&User{}, u, nil, id); err != nil {
		logging.Err.Errorf("update user err: %+v\n", err)
		return err
	}
	if err := addRoles(u); err != nil {
		return err
	}
	return nil
}

// addRoles add roles for user
func addRoles(user *User) error {
	if len(user.RoleIds) > 0 {
		userId := strconv.FormatUint(uint64(user.ID), 10)
		if _, err := easygorm.Egm.Enforcer.DeleteRolesForUser(userId); err != nil {
			logging.Err.Errorf("add role to user,del role  err: %+v\n", err)
			return err
		}

		for _, roleId := range user.RoleIds {
			roleId := strconv.FormatUint(uint64(roleId), 10)
			if _, err := easygorm.Egm.Enforcer.AddRoleForUser(userId, roleId); err != nil {
				logging.Err.Errorf("add role to user,add role err: %+v\n", err)
				return err
			}
		}
	}
	return nil
}
