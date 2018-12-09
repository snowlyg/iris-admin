package main

import (
	"IrisApiProject/config"
	"IrisApiProject/controllers"
	"IrisApiProject/database"
	"IrisApiProject/middleware"
	"IrisApiProject/models"
	"github.com/betacraft/yaag/yaag"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
)

/**
 * 初始化 iris app
 * @method NewApp
 * @return  {[type]}  api      *iris.Application  [iris app]
 */
func newApp() (api *iris.Application) {

	api = iris.New()
	api.Use(logger.New())

	api.OnErrorCode(iris.StatusNotFound, controllers.NotFound)
	api.OnErrorCode(iris.StatusInternalServerError, controllers.InternalServerError)

	//同步模型数据表
	database.DB.AutoMigrate(new(models.Users))

	iris.RegisterOnInterrupt(func() {
		database.DB.Close()
	})

	//"github.com/iris-contrib/middleware/cors"
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	appName := config.Conf.Get("app.name").(string)
	appDoc := config.Conf.Get("app.doc").(string)
	appUrl := config.Conf.Get("app.url").(string)

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: appName,
		DocPath:  appDoc + "/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": appUrl,
			"Staging":    "",
		},
	})

	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Use(middleware.NewYaag()) // <- IMPORTANT, register the middleware.
		v1.Post("/admin/login", controllers.CUserLogin)
		v1.PartyFunc("/admin", func(admin router.Party) {
			admin.Use(middleware.JwtHandler().Serve, middleware.AuthToken)
			admin.Get("/logout", controllers.CUserLogout)

			admin.PartyFunc("/users", func(users router.Party) {
				users.Get("/", controllers.CGetAllUsers)
				users.Get("/{id:uint}", controllers.CGetUser)
				users.Post("/", controllers.CCreateUser)
				users.Post("/{id:uint}/update", controllers.CUpdateUser)
				users.Get("/{id:uint}/frozen", controllers.CFrozenUser)
				users.Get("/{id:uint}/audit", controllers.CSetUserAudit)
				users.Get("/{id:uint}/refrozen", controllers.CRefrozenUser)
				users.Delete("/{id:uint}", controllers.CDeleteUser)
				users.Get("/profile", controllers.CGetProfile)
			})
		})
	}

	return
}

func main() {

	app := newApp()

	addr := config.Conf.Get("app.addr").(string)
	app.Run(iris.Addr(addr))
}
