package duser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/snowlyg/blog/application/libs"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/models"
)

const ModelName = "用户管理"

type UserResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Intro     string `json:"introduction"`
	Avatar    string `json:"avatar"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

type UserReq struct {
	Name     string `json:"name" `
	Username string `json:"username"`
	Password string `json:"password"`
	Intro    string `json:"introduction"`
	Avatar   string `json:"avatar"`
}

func (u *UserResponse) ModelName() string {
	return ModelName
}

func (u *UserResponse) Model() *models.User {
	return &models.User{}
}

func (u *UserResponse) All(name, sort, orderBy string, page, pageSize int) (map[string]interface{}, error) {
	var count int64
	var users []*UserResponse
	db := easygorm.GetEasyGormDb().Model(u.Model())
	if len(name) > 0 {
		db = db.Where("name", "like", fmt.Sprintf("%%%s%%", name))
	}
	err := db.Count(&count).Error
	if err != nil {
		logging.ErrorLogger.Errorf("get list count err ", err)
		return nil, err
	}

	err = db.Scopes(easygorm.PaginateScope(page, pageSize, sort, orderBy)).Find(&users).Error
	if err != nil {
		logging.ErrorLogger.Errorf("get list data err ", err)
		return nil, err
	}

	list := map[string]interface{}{"items": users, "total": count, "limit": pageSize}
	return list, nil
}

func (u *UserResponse) FindByUserName(username string) error {
	err := easygorm.GetEasyGormDb().Model(u.Model()).Where("username = ?", username).Find(u).Error
	if err != nil {
		logging.ErrorLogger.Errorf("find user by username ", username, " err ", err)
		return err
	}
	return nil
}

func (u *UserResponse) Create(object map[string]interface{}) error {
	if username, ok := object["Username"].(string); ok {
		err := u.FindByUserName(username)
		if err != nil {
			logging.ErrorLogger.Errorf("create user find by username get err ", err)
			return err
		}

		if u.Id > 0 {
			return errors.New(fmt.Sprintf("username %s is being used", username))
		}
	}

	err := easygorm.GetEasyGormDb().Model(u.Model()).Create(object).Error
	if err != nil {
		logging.ErrorLogger.Errorf("create data err ", err)
		return err
	}

	return nil
}

func (u *UserResponse) Update(id uint, object map[string]interface{}) error {
	err := u.Find(id)
	if err != nil {
		return err
	}
	if u.Username == "username" {
		return errors.New("不能编辑管理员")
	}
	if username, ok := object["Username"].(string); ok {
		err := u.FindByUserName(username)
		if err != nil {
			logging.ErrorLogger.Errorf("create user find by username get err ", err)
			return err
		}

		if u.Id > 0 && u.Id != id {
			return errors.New(fmt.Sprintf("username %s is being used", username))
		}
	}
	err = easygorm.GetEasyGormDb().Model(u.Model()).Where("id = ?", id).Updates(object).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserResponse) Find(id uint) error {
	err := easygorm.GetEasyGormDb().Model(u.Model()).Where("id = ?", id).Find(u).Error
	if err != nil {
		logging.ErrorLogger.Errorf("find user err ", err)
		return err
	}
	return nil
}

func (u *UserResponse) Delete(id uint) error {
	err := easygorm.GetEasyGormDb().Delete(u.Model(), id).Error
	if err != nil {
		logging.ErrorLogger.Errorf("delete user by id get  err ", err)
		return err
	}
	return nil
}

// AddRoleForUser add roles for user
func AddRoleForUser(user *models.User) error {
	if len(user.RoleIds) == 0 {
		return nil
	}

	var err error
	var roleIds []string
	var oldRoleIds []string

	userId := strconv.FormatUint(uint64(user.ID), 10)
	oldRoleIds, err = easygorm.GetEasyGormEnforcer().GetRolesForUser(userId)
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

	if _, err := easygorm.GetEasyGormEnforcer().AddRolesForUser(userId, roleIds); err != nil {
		logging.ErrorLogger.Errorf("add role to user role failed: %+v\n", err)
		return err
	}

	return nil
}

func (u *UserResponse) Profile(id uint) error {
	return u.Find(id)
}
