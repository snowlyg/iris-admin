package models

/**
*初始化系统 账号 权限 角色
 */
func CreateSystemData() {
	perm := CreateSystemAdminPermission() //初始化权限
	permIds := []uint{perm.ID}
	role := CreateSystemAdminRole(permIds) //初始化角色
	if role.ID != 0 {
		CreateSystemAdmin(role.ID) //初始化管理员
	}
}
