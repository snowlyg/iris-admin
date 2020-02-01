package models

import (
	"errors"
	"fmt"

	"IrisAdminApi/database"
	"IrisAdminApi/transformer"
	"IrisAdminApi/validates"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

/**
*初始化系统 账号 权限 角色
 */
func CreateSystemData(rc *transformer.Conf, perms []*validates.PermissionRequest) {
	if rc.App.CreateSysData {
		permIds := CreateSystemAdminPermission(perms) //初始化权限
		role := CreateSystemAdminRole(permIds)        //初始化角色
		if role.ID != 0 {
			CreateSystemAdmin(role.ID, rc) //初始化管理员
		}
	}
}

/**
*创建系统管理员
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemAdmin(roleId uint, rc *transformer.Conf) {
	aul := &validates.CreateUpdateUserRequest{
		Username: rc.TestData.UserName,
		Password: rc.TestData.Pwd,
		Name:     rc.TestData.Name,
		RoleIds:  []uint{roleId},
	}

	user := NewUserByStruct(aul)
	user.GetUserByUsername()
	if user.ID == 0 {
		user.CreateUser(aul)
	}
}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreateSystemAdminRole(permIds []uint) *Role {
	rr := &validates.RoleRequest{
		Name:        "admin",
		DisplayName: "管理员",
		Description: "管理员",
	}
	role := NewRoleByStruct(rr)
	role.GetRoleByName()
	if role.ID == 0 {
		role.CreateRole(permIds)
	}

	return role
}

/**
 * 创建系统权限
 * @return
 */
func CreateSystemAdminPermission(perms []*validates.PermissionRequest) []uint {
	var permIds []uint
	for _, perm := range perms {
		p := NewPermission(0, perm.Name, perm.Act)
		p.GetPermissionByNameAct()
		if p.ID != 0 {
			continue
		}
		p.CreatePermission()
		permIds = append(permIds, p.ID)
	}
	return permIds
}

func IsNotFound(err error) {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); !ok && err != nil {
		color.Red(fmt.Sprintf("error :%v \n ", err))
	}
}

/**
 * 获取列表
 * @method MGetAll
 * @param  {[type]} string string    [description]
 * @param  {[type]} orderBy string    [description]
 * @param  {[type]} relation string    [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAll(string, orderBy string, offset, limit int) *gorm.DB {
	db := database.GetGdb()
	if len(orderBy) > 0 {
		db.Order(orderBy + "desc")
	} else {
		db.Order("created_at desc")
	}
	if len(string) > 0 {
		db.Where("name LIKE  ?", "%"+string+"%")
	}
	if offset > 0 {
		db.Offset((offset - 1) * limit)
	}
	if limit > 0 {
		db.Limit(limit)
	}
	return db
}

func DelAllData() {
	database.GetGdb().Unscoped().Delete(&OauthToken{})
	database.GetGdb().Unscoped().Delete(&Permission{})
	database.GetGdb().Unscoped().Delete(&Role{})
	database.GetGdb().Unscoped().Delete(&User{})
	database.GetGdb().Exec("DELETE FROM casbin_rule;")
}
