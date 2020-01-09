package routes

import (
	"IrisAdminApi/controllers"
	"IrisAdminApi/middleware"
	"github.com/betacraft/yaag/irisyaag"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

func Register(api *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	api.Post("/v1/admin/login", controllers.UserLogin)
	main := api.Party("/", crs).AllowMethods(iris.MethodOptions)
	{
		home := main.Party("/")
		home.Get("/", func(ctx iris.Context) { // 首页模块
			_ = ctx.View("index.html")
		})

		v1 := main.Party("/v1")
		{
			v1.PartyFunc("/admin", func(admin iris.Party) {
				admin.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.
				admin.Use(middleware.JwtHandler().Serve, middleware.AuthToken)
				admin.Get("/logout", controllers.UserLogout)

				admin.PartyFunc("/users", func(users iris.Party) {
					users.Get("/", controllers.GetAllUsers)
					users.Get("/{id:uint}", controllers.GetUser)
					users.Post("/", controllers.CreateUser)
					users.Put("/{id:uint}", controllers.UpdateUser)
					users.Delete("/{id:uint}", controllers.DeleteUser)
					users.Get("/profile", controllers.GetProfile)
				})
				admin.PartyFunc("/roles", func(roles iris.Party) {
					roles.Get("/", controllers.GetAllRoles)
					roles.Get("/{id:uint}", controllers.GetRole)
					roles.Post("/", controllers.CreateRole)
					roles.Put("/{id:uint}", controllers.UpdateRole)
					roles.Delete("/{id:uint}", controllers.DeleteRole)
				})
				admin.PartyFunc("/permissions", func(permissions iris.Party) {
					permissions.Get("/", controllers.GetAllPermissions)
					permissions.Get("/{id:uint}", controllers.GetPermission)
					permissions.Post("/", controllers.CreatePermission)
					permissions.Post("/import", controllers.ImportPermission)
					permissions.Put("/{id:uint}", controllers.UpdatePermission)
					permissions.Delete("/{id:uint}", controllers.DeletePermission)
				})
			})
		}
	}

	api.Get("/apiDoc", func(ctx iris.Context) { // 首页模块
		_ = ctx.View("apiDoc/index.html")
	})

}
