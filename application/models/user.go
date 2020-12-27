package models

import (
	"github.com/snowlyg/blog/application/libs"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"gorm.io/gorm"
	"strconv"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null; type:varchar(60)" json:"name" `
	Username string `gorm:"unique;not null;type:varchar(60)" json:"username"`
	Password string `gorm:"type:varchar(100)" json:"password"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"introduction"`
	Avatar   string `gorm:"type:longText" json:"avatar"`
	RoleIds  []uint `gorm:"-" json:"role_ids"`
}

//// GetUser get user
//func GetUser(search *easygorm.Search) (*User, error) {
//	user := &User{}
//	if err := easygorm.First(user, search); err != nil {
//		logging.ErrorLogger.Errorf("get user err: %+v\n", err)
//		return user, err
//	}
//	return user, nil
//}
//
//// DeleteUser del user . if user's username is username ,can't del it.
//func DeleteUser(id uint) error {
//	user := &User{}
//	if err := easygorm.FindById(user, id); err != nil {
//		return err
//	}
//	if user.Username == "username" {
//		err := errors.New(fmt.Sprintf("不能删除管理员 : %s \n ", user.Username))
//		return err
//	}
//	if err := easygorm.EasyGorm.DB.Delete(user, id).Error; err != nil {
//		logging.ErrorLogger.Errorf("del user err: %+v\n", err)
//		return err
//	}
//	return nil
//}
//
//// GetAllUsers get all users
//func GetAllUsers(s *easygorm.Search) ([]*User, int64, error) {
//	var users []*User
//	count, err := easygorm.Paginate(&User{}, &users, s)
//	if err != nil {
//		logging.ErrorLogger.Errorf("get all users err: %+v\n", err)
//		return nil, count, err
//	}
//	return users, count, nil
//}
//
//// CreateUser create user
//func (u *User) CreateUser() error {
//	u.Password = libs.HashPassword(u.Password)
//	if err := easygorm.Create(u); err != nil {
//		logging.ErrorLogger.Errorf("create user err: %+v\n", err)
//		return err
//	}
//	if err := addRoles(u); err != nil {
//		return err
//	}
//	return nil
//}
//
//// UpdateUserById update user by id
//func UpdateUserById(id uint, u *User) error {
//	if len(u.Password) > 0 {
//		u.Password = libs.HashPassword(u.Password)
//	}
//	if err := easygorm.Update(&User{}, u, nil, id); err != nil {
//		logging.ErrorLogger.Errorf("update user err: %+v\n", err)
//		return err
//	}
//	if err := addRoles(u); err != nil {
//		return err
//	}
//	return nil
//}

// AddRoleForUser add roles for user
func AddRoleForUser(user *User) error {
	if len(user.RoleIds) == 0 {
		return nil
	}

	var err error
	var roleIds []string
	var oldRoleIds []string

	userId := strconv.FormatUint(uint64(user.ID), 10)
	oldRoleIds, err = easygorm.EasyGorm.Enforcer.GetRolesForUser(userId)
	if err != nil {
		logging.ErrorLogger.Errorf("add role to user,del role  err: %+v\n", err)
		return err
	}

	for _, roleId := range user.RoleIds {
		roleId := strconv.FormatUint(uint64(roleId), 10)
		if len(oldRoleIds) > 0 && libs.InArrayS(oldRoleIds, roleId) {
			continue
		}

		roleIds = append(roleIds, roleId)
	}

	if _, err := easygorm.EasyGorm.Enforcer.AddRolesForUser(userId, roleIds); err != nil {
		logging.ErrorLogger.Errorf("add role to user role failed: %+v\n", err)
		return err
	}

	return nil
}
