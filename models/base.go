package models

/**
*初始化系统 账号 权限 角色
 */
func CreateSystemData() {
	CreateSystemAdminPermission()   //初始化权限
	role := CreateSystemAdminRole() //初始化角色

	if role.ID != 0 {
		CreateSystemAdmin(role.ID) //初始化管理员
	}
}
