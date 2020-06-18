package routes

import (
	"github.com/betacraft/yaag/irisyaag"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/backend/config"
	"github.com/snowlyg/IrisAdminApi/backend/controllers"
	"github.com/snowlyg/IrisAdminApi/backend/middleware"
	"github.com/snowlyg/IrisAdminApi/backend/sysinit"
)

func App(api *iris.Application) {
	//api.Favicon("./static/favicons/favicon.ico")
	app := api.Party("/", middleware.CrsAuth()).AllowMethods(iris.MethodOptions)
	{
		app.HandleDir("/static", config.Root+"resources/app/static")
		app.HandleDir("/record", config.Config.RecordPath) // 视频记录地址
		app.Get("/", func(ctx iris.Context) {              // 首页模块
			_ = ctx.View("app/index.html")
		})

		v1 := app.Party("/v1")
		{
			v1.Post("/admin/login", controllers.UserLogin)
			v1.Use(irisyaag.New())
			v1.PartyFunc("/admin", func(app iris.Party) {
				app.Get("/resetData", controllers.ResetData)
				casbinMiddleware := middleware.New(sysinit.Enforcer)               //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
				app.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
				app.Get("/logout", controllers.UserLogout).Name = "退出"

				app.PartyFunc("/users", func(users iris.Party) {
					users.Get("/", controllers.GetAllUsers).Name = "用户列表"
					users.Get("/{id:uint}", controllers.GetUser).Name = "用户详情"
					users.Post("/", controllers.CreateUser).Name = "创建用户"
					users.Put("/{id:uint}", controllers.UpdateUser).Name = "编辑用户"
					users.Delete("/{id:uint}", controllers.DeleteUser).Name = "删除用户"
					users.Get("/profile", controllers.GetProfile).Name = "个人信息"
				})
				app.PartyFunc("/roles", func(roles iris.Party) {
					roles.Get("/", controllers.GetAllRoles).Name = "角色列表"
					roles.Get("/{id:uint}", controllers.GetRole).Name = "角色详情"
					roles.Post("/", controllers.CreateRole).Name = "创建角色"
					roles.Put("/{id:uint}", controllers.UpdateRole).Name = "编辑角色"
					roles.Delete("/{id:uint}", controllers.DeleteRole).Name = "删除角色"
				})
				app.PartyFunc("/permissions", func(permissions iris.Party) {
					permissions.Get("/", controllers.GetAllPermissions).Name = "权限列表"
					permissions.Get("/{id:uint}", controllers.GetPermission).Name = "权限详情"
					permissions.Post("/import", controllers.ImportPermission).Name = "导入权限"
					permissions.Post("/", controllers.CreatePermission).Name = "创建权限"
					permissions.Put("/{id:uint}", controllers.UpdatePermission).Name = "编辑权限"
					permissions.Delete("/{id:uint}", controllers.DeletePermission).Name = "删除权限"
				})
			})
		}
	}

}
