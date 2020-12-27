package models

import (
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"gorm.io/gorm"
	"strconv"
)

type Role struct {
	gorm.Model

	Name        string     `gorm:"unique;not null; type:varchar(256)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string     `gorm:"type:varchar(256)" json:"display_name" comment:"显示名称"`
	Description string     `gorm:"type:varchar(256)" json:"description" comment:"描述"`
	Perms       [][]string `gorm:"-" json:"perms" comment:"权限 name, act "`
}

//
//// GetRoleById get role by it
//func GetRoleById(id uint) (*Role, error) {
//	t := &Role{}
//	err := easygorm.FindById(&Role{}, id)
//	if err != nil {
//		logging.ErrorLogger.Errorf("get role by id err: %+v", err)
//		return t, err
//	}
//
//	return t, nil
//}
//
//// GetRole get role
//func GetRole(s *easygorm.Search) (*Role, error) {
//	t := &Role{}
//	err := easygorm.First(t, s)
//	if err != nil {
//		logging.ErrorLogger.Errorf("get role err: %+v", err)
//		return t, err
//	}
//
//	return t, nil
//}
//
//// DeleteRoleById del role by id
//func DeleteRoleById(id uint) error {
//	r := &Role{}
//	err := easygorm.DeleteById(r, id)
//	if err != nil {
//		logging.ErrorLogger.Errorf("del role by id err: %+v", err)
//		return err
//	}
//
//	return nil
//}
//
//// GetAllRoles get all roles
//func GetAllRoles(s *easygorm.Search) ([]*Role, int64, error) {
//	var roles []*Role
//	count, err := easygorm.Paginate(&Role{}, &roles, s)
//	if err != nil {
//		logging.ErrorLogger.Errorf("get all role err: %+v", err)
//		return nil, count, err
//	}
//
//	return roles, count, nil
//}
//
//// CreateRole create role
//func (r *Role) CreateRole() error {
//	if err := easygorm.Create(r); err != nil {
//		logging.ErrorLogger.Errorf("create role err: %+v", err)
//		return err
//	}
//	if err := addPerms(r.PermIds, r); err != nil {
//		return err
//	}
//	return nil
//}

// AddPermForRole add perms
func AddPermForRole(role *Role) error {
	if len(role.Perms) == 0 {
		logging.DebugLogger.Debugf("no perms")
		return nil
	}

	var newPerms [][]string
	roleId := strconv.FormatUint(uint64(role.ID), 10)
	oldPerms := easygorm.EasyGorm.Enforcer.GetPermissionsForUser(roleId)

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
	_, err = easygorm.EasyGorm.Enforcer.AddPolicies(newPerms)
	if err != nil {
		logging.ErrorLogger.Errorf("add policy err: %+v", err)
		return err
	}

	return nil
}

//
//// UpdateRole update role
//func UpdateRole(id uint, r *Role) error {
//	if err := easygorm.Update(&Role{}, r, []interface{}{"DisplayName", "Description"}, id); err != nil {
//		logging.ErrorLogger.Errorf("update role err: %+v", err)
//		return err
//	}
//	if err := addPerms(r.PermIds, r); err != nil {
//		return err
//	}
//	return nil
//}
//
//// RolePermissions get role's permissions
//func (r *Role) RolePermissions() []*Permission {
//	perms := easygorm.GetPermissionsForUser(r.ID)
//	var ps []*Permission
//	for _, perm := range perms {
//		if len(perm) >= 3 && len(perm[1]) > 0 && len(perm[2]) > 0 {
//			s := &easygorm.Search{
//				Fields: []*easygorm.Field{
//					{
//						Key:       "name",
//						Condition: "=",
//						Value:     perm[1],
//					},
//					{
//						Key:       "act",
//						Condition: "=",
//						Value:     perm[2],
//					},
//				},
//			}
//			p, err := GetPermission(s)
//			if err == nil && p.ID > 0 {
//				ps = append(ps, p)
//			}
//		}
//	}
//	return ps
//}
