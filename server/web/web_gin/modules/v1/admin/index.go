package admin

import (
	"github.com/gin-gonic/gin"
)

func Group(group *gin.RouterGroup) {
	adminRouter := group.Group("/admin")
	{
		adminRouter.GET("/", GetAll)
		adminRouter.GET("/{id:uint}", GetAdmin)
		adminRouter.POST("/", CreateAdmin)
		adminRouter.POST("/{id:uint}", UpdateAdmin)
		adminRouter.DELETE("/{id:uint}", DeleteAdmin)

		adminRouter.GET("/profile", Profile)
		adminRouter.POST("/change_avatar", ChangeAvatar)
	}
}
