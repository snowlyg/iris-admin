package drole

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/snowlyg/iris-admin/application/libs/easygorm"
	"github.com/snowlyg/iris-admin/application/libs/logging"
	"github.com/snowlyg/iris-admin/application/models"
)

const modelName = "角色管理"
const adminRoleName = "admin"

type RoleResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

type RoleReq struct {
	Name        string `json:"name" `
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

func (r *RoleResponse) ModelName() string {
	return modelName
}

func Model() *models.Role {
	return &models.Role{}
}

func (r *RoleResponse) All(name, sort, orderBy string, page, pageSize int) (map[string]interface{}, error) {
	var count int64
	var roles []*RoleResponse
	db := easygorm.GetEasyGormDb().Model(Model())
	if len(name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	err := db.Count(&count).Error
	if err != nil {
		logging.ErrorLogger.Errorf("get list count err ", err)
		return nil, err
	}
	err = db.Scopes(easygorm.PaginateScope(page, pageSize, sort, orderBy)).Find(&roles).Error
	if err != nil {
		logging.ErrorLogger.Errorf("get list data err ", err)
		return nil, err
	}
	list := map[string]interface{}{"items": roles, "total": count, "limit": pageSize}
	return list, nil
}

func (r *RoleResponse) FindByName(name string) error {
	err := easygorm.GetEasyGormDb().Model(Model()).Where("name = ?", name).First(r).Error
	if err != nil {
		logging.ErrorLogger.Errorf("find role by name get err ", err)
		return err
	}
	return nil
}

func (r *RoleResponse) Create(object map[string]interface{}) error {
	if name, ok := object["Name"].(string); ok {
		err := r.FindByName(name)
		if err != nil {
			logging.ErrorLogger.Errorf("create role find by name get err ", err)
			return err
		}

		if r.Id > 0 {
			return errors.New(fmt.Sprintf("name %s is being used", name))
		}
	}

	err := easygorm.GetEasyGormDb().Model(Model()).Create(object).Error
	if err != nil {
		logging.ErrorLogger.Errorf("create data err ", err)
		return err
	}

	return nil
}
func (r *RoleResponse) Update(id uint, object map[string]interface{}) error {
	err := r.First(id)
	if err != nil {
		return err
	}
	if r.Name == adminRoleName {
		return errors.New("不能编辑管理员角色")
	}
	if name, ok := object["Name"].(string); ok {
		err := r.FindByName(name)
		if err != nil {
			logging.ErrorLogger.Errorf("create role find by name get err ", err)
			return err
		}

		if r.Id > 0 && r.Id != id {
			return errors.New(fmt.Sprintf("name %s is being used", name))
		}
	}
	err = easygorm.GetEasyGormDb().Model(Model()).Where("id = ?", id).Updates(object).Error
	if err != nil {
		logging.ErrorLogger.Errorf("update role  get err ", err)
		return err
	}
	return nil
}

func (r *RoleResponse) First(id uint) error {
	err := easygorm.GetEasyGormDb().Model(Model()).Where("id = ?", id).First(r).Error
	if err != nil {
		logging.ErrorLogger.Errorf("find role by id get  err ", err)
		return err
	}
	return nil
}

func (r *RoleResponse) Delete(id uint) error {
	err := easygorm.GetEasyGormDb().Unscoped().Delete(Model(), id).Error
	if err != nil {
		logging.ErrorLogger.Errorf("delete role by id get  err ", err)
		return err
	}
	return nil
}

func FindInId(ids []string) ([]*RoleResponse, error) {
	var roles []*RoleResponse
	err := easygorm.GetEasyGormDb().Model(Model()).Where("id in ?", ids).Find(&roles).Error
	if err != nil {
		logging.ErrorLogger.Errorf("find role by id get  err ", err)
		return nil, err
	}
	return roles, nil
}

// AddPermForRole add perms
func AddPermForRole(role *models.Role) error {
	if len(role.Perms) == 0 {
		logging.DebugLogger.Debugf("no perms")
		return nil
	}

	var newPerms [][]string
	roleId := strconv.FormatUint(uint64(role.ID), 10)
	oldPerms := easygorm.GetEasyGormEnforcer().GetPermissionsForUser(roleId)

	for _, perm := range role.Perms {
		var in bool
		for _, oldPerm := range oldPerms {
			if roleId == oldPerm[0] && perm[0] == oldPerm[1] && perm[1] == oldPerm[2] {
				in = true
				continue
			}
		}

		if !in {
			newPerms = append(newPerms, append([]string{roleId}, perm...))
		}
	}

	logging.DebugLogger.Debugf("new perms", newPerms)

	var err error
	_, err = easygorm.GetEasyGormEnforcer().AddPolicies(newPerms)
	if err != nil {
		logging.ErrorLogger.Errorf("add policy err: %+v", err)
		return err
	}

	return nil
}
