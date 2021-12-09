package authority

import "github.com/gin-gonic/gin"

func Group(group *gin.RouterGroup) {

	authRouter := group.Group("/authority")
	{
		authRouter.POST("/createAuthority", CreateAuthority)                 // 创建角色
		authRouter.PUT("/updateAuthority", UpdateAuthority)                  // 更新角色
		authRouter.POST("/copyAuthority", CopyAuthority)                     // 更新角色
		authRouter.POST("/getAuthorityList", GetAuthorityList)               // 获取角色列表
		authRouter.POST("/getAdminAuthorityList", GetAdminAuthorityList)     // 获取员工角色列表
		authRouter.POST("/getTenancyAuthorityList", GetTenancyAuthorityList) // 获取商户角色列表
		authRouter.POST("/getGeneralAuthorityList", GetGeneralAuthorityList) // 获取普通用户角色列表
		authRouter.DELETE("/deleteAuthority", DeleteAuthority)               // 删除角色
	}
}
