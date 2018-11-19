package main

import (
	"IrisYouQiKangApi/controllers"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
)

func newApp() *iris.Application {
	api := iris.New()
	api.Use(logger.New())
	api.OnErrorCode(iris.StatusNotFound, controllers.NotFound)
	api.OnErrorCode(iris.StatusInternalServerError, controllers.InternalServerError)

	jwtHandler := jwtHandler()

	// or	"github.com/iris-contrib/middleware/cors"
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Post("/admin/login", controllers.UserAdminLogin)

		v1.PartyFunc("/admin", func(admin router.Party) {
			admin.Use(jwtHandler.Serve)
			admin.Get("/",controllers.GetHomeData)
			admin.PartyFunc("/users", func(users router.Party) {
				users.Get("/profile", controllers.GetProfile)
			})
		})

	}

	return api
}
