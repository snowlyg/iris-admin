package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/controllers"
	"github.com/snowlyg/IrisAdminApi/middleware"
	"github.com/snowlyg/IrisAdminApi/sysinit"
)

const maxSize = 5 << 20 // 5MB

func App(api *iris.Application) {
	api.UseRouter(middleware.CrsAuth())
	app := api.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.HandleDir("/uploads", iris.Dir("./uploads"))
		if config.Config.Bindata {
			app.Get("/", func(ctx iris.Context) { // 首页模块
				_ = ctx.View("index.html")
			})
		}

		v1 := app.Party("/v1")
		{
			v1.PartyFunc("/article", func(aritcle iris.Party) {
				aritcle.Get("/", controllers.GetAllArticles)
				aritcle.Get("/{id:uint}", controllers.GetArticle)
			})
			v1.Post("/admin/login", controllers.UserLogin)
			v1.PartyFunc("/admin", func(admin iris.Party) {
				casbinMiddleware := middleware.New(sysinit.Enforcer)                 //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
				admin.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
				admin.Post("/logout", controllers.UserLogout).Name = "退出"
				admin.Get("/profile", controllers.GetProfile).Name = "个人信息"
				admin.Post("/upload_file", iris.LimitRequestBodySize(maxSize+1<<20), controllers.UploadFile).Name = "上传文件"
				admin.PartyFunc("/article", func(aritcle iris.Party) {
					aritcle.Get("/", controllers.GetAllArticles).Name = "文章列表"
					aritcle.Get("/{id:uint}", controllers.GetArticle).Name = "文章详情"
					aritcle.Post("/", controllers.CreateArticle).Name = "创建文章"
					aritcle.Put("/{id:uint}", controllers.UpdateArticle).Name = "编辑文章"
					aritcle.Delete("/{id:uint}", controllers.DeleteArticle).Name = "删除文章"
				})

				admin.PartyFunc("/users", func(users iris.Party) {
					users.Get("/", controllers.GetAllUsers).Name = "用户列表"
					users.Get("/{id:uint}", controllers.GetUser).Name = "用户详情"
					users.Post("/", controllers.CreateUser).Name = "创建用户"
					users.Put("/{id:uint}", controllers.UpdateUser).Name = "编辑用户"
					users.Delete("/{id:uint}", controllers.DeleteUser).Name = "删除用户"
				})
				admin.PartyFunc("/roles", func(roles iris.Party) {
					roles.Get("/", controllers.GetAllRoles).Name = "角色列表"
					roles.Get("/{id:uint}", controllers.GetRole).Name = "角色详情"
					roles.Post("/", controllers.CreateRole).Name = "创建角色"
					roles.Put("/{id:uint}", controllers.UpdateRole).Name = "编辑角色"
					roles.Delete("/{id:uint}", controllers.DeleteRole).Name = "删除角色"
				})
				admin.PartyFunc("/permissions", func(permissions iris.Party) {
					permissions.Get("/", controllers.GetAllPermissions).Name = "权限列表"
					permissions.Get("/{id:uint}", controllers.GetPermission).Name = "权限详情"
					permissions.Post("/", controllers.CreatePermission).Name = "创建权限"
					permissions.Put("/{id:uint}", controllers.UpdatePermission).Name = "编辑权限"
					permissions.Delete("/{id:uint}", controllers.DeletePermission).Name = "删除权限"
				})
			})
		}
	}

}
