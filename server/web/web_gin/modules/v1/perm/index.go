package perm

import (
	"github.com/gin-gonic/gin"
)

func Group(group *gin.RouterGroup) {
	ApiRouter := group.Group("/perms")
	{
		ApiRouter.GET("/getApiList", GetApiList)              // 获取Api列表
		ApiRouter.GET("/getAllApis", GetAllApis)              // 获取所有api
		ApiRouter.GET("/getApiById/:id", GetApiById)          // 获取单条Api消息
		ApiRouter.POST("/createApi", CreateApi)               // 创建Api
		ApiRouter.DELETE("/deleteApi", DeleteApi)             // 删除Api
		ApiRouter.PUT("/updateApi/:id", UpdateApi)            // 更新api
		ApiRouter.DELETE("/deleteApisByIds", DeleteApisByIds) // 删除选中api
	}
}
