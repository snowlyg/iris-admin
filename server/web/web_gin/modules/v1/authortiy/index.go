package authority

import "github.com/gin-gonic/gin"

func Group(group *gin.RouterGroup) {

	ApiRouter := group.Group("/authority")
	{
		ApiRouter.POST("/createAuthority", CreateAuthority)                 // 创建角色
		ApiRouter.PUT("/updateAuthority", UpdateAuthority)                  // 更新角色
		ApiRouter.POST("/copyAuthority", CopyAuthority)                     // 更新角色
		ApiRouter.POST("/getAuthorityList", GetAuthorityList)               // 获取角色列表
		ApiRouter.POST("/getAdminAuthorityList", GetAdminAuthorityList)     // 获取员工角色列表
		ApiRouter.POST("/getTenancyAuthorityList", GetTenancyAuthorityList) // 获取商户角色列表
		ApiRouter.POST("/getGeneralAuthorityList", GetGeneralAuthorityList) // 获取普通用户角色列表
		ApiRouter.DELETE("/deleteAuthority", DeleteAuthority)               // 删除角色
	}
}
