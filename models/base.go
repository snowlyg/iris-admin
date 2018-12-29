package models

/**
*初始化系统 账号 权限 角色
 */
func CreaterSystemData() {
	CreaterSystemAdmin()     //初始化管理员
	CreaterSystemAdminRole() //初始化角色
}
