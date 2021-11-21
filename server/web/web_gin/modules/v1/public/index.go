package public

import (
	"github.com/gin-gonic/gin"
)

// Group 认证模块
func Group(group *gin.RouterGroup) {
	publicRouter := group.Group("/public")
	{
		publicRouter.GET("/captcha", Captcha)
		publicRouter.POST("/admin/login", AdminLogin)
		// publicRouter.POST("/merchant/login", ClientLogin)
		// publicRouter.POST("/device/login", LoginDevice)
		// publicRouter.GET("/mini/auth", AuthMini)
		// publicRouter.POST("/mini/login", LoginMini)
		// publicRouter.POST("/mini/getPhoneNumber", GetPhoneNumber)
		// publicRouter.GET("/region/:code", Region)
		// publicRouter.GET("/getRegionList", RegionList)
		// publicRouter.GET("/getRefundMessage", GetRefundMessage)
	}
	group.GET("/logout", Logout) // 退出
	group.GET("/clean", Clear)   //清空授权
}
