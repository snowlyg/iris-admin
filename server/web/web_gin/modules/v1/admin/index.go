package admin

import (
	"github.com/gin-gonic/gin"
)

func Group(group *gin.RouterGroup) {
	ApiRouter := group.Group("/perms")
	{
		ApiRouter.GET("/", GetAll)
		ApiRouter.GET("/{id:uint}", GetAdmin)
		ApiRouter.POST("/", CreateAdmin)
		ApiRouter.POST("/{id:uint}", UpdateAdmin)
		ApiRouter.DELETE("/{id:uint}", DeleteAdmin)
		ApiRouter.GET("/logout", Logout)
		ApiRouter.GET("/clear", Clear)
		ApiRouter.GET("/profile", Profile)
		ApiRouter.POST("/change_avatar", ChangeAvatar)
	}
}
