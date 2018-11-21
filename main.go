package main

import (
	"IrisYouQiKangApi/controllers"
	"IrisYouQiKangApi/middleware"
	"IrisYouQiKangApi/models"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
)

func main() {

	api := iris.New()
	api.Use(logger.New())

	api.OnErrorCode(iris.StatusNotFound, controllers.NotFound)
	api.OnErrorCode(iris.StatusInternalServerError, controllers.InternalServerError)

	iris.RegisterOnInterrupt(func() {
		models.DB.Close()
	})

	// or	"github.com/iris-contrib/middleware/cors"
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	{

		v1.Post("/admin/login", controllers.UserAdminLogin)

		v1.PartyFunc("/admin", func(admin router.Party) {
			admin.Use(middleware.JwtHandler().Serve, middleware.AuthToken)
			admin.Get("/", controllers.GetHomeData)
			admin.Get("/logout", controllers.UserAdminLogout)
			admin.PartyFunc("/users", func(users router.Party) {
				users.Get("/", controllers.GetAllUsers)
				users.Get("/profile", controllers.GetProfile)
			})
		})

	}

	api.Run(iris.Addr(":80"))
}
