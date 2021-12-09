package public

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
)

// Group 认证模块
func Group(group *gin.RouterGroup) {
	group.GET("/public/captcha", Captcha)
	group.POST("/public/admin/login", AdminLogin)
	// group.POST("/public/merchant/login", ClientLogin)
	// group.POST("/public/device/login", LoginDevice)
	// group.GET("/public/mini/auth", AuthMini)
	// group.POST("/public/mini/login", LoginMini)
	// group.POST("/public/mini/getPhoneNumber", GetPhoneNumber)
	// group.GET("/public/region/:code", Region)
	// group.GET("/public/getRegionList", RegionList)
	// group.GET("/public/getRefundMessage", GetRefundMessage)
	publicRouter := group.Group("/public")
	{
		publicRouter.Use(middleware.Auth(), middleware.CasbinHandler(), middleware.OperationRecord())
		publicRouter.GET("/logout", Logout) // 退出
		publicRouter.GET("/clean", Clear)   //清空授权
	}

}
