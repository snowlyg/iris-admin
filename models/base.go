package models

/**
*初始化系统 账号 权限 角色
 */
func CreateSystemData() {
	CreateSystemAdmin()           //初始化管理员
	CreateSystemAdminRole()       //初始化角色
	CreateSystemAdminPermission() //初始化权限
}
