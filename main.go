package main

import (
	"IrisYouQiKangApi/controllers"
	"IrisYouQiKangApi/middleware"
	"IrisYouQiKangApi/models"
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
)

func init() {

}

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

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: models.Config.App.Name,
		DocPath:  models.Config.App.Doc + "/index.html",
		BaseUrls: map[string]string{"Production": models.Config.App.Url, "Staging": ""},
	})
	api.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.

	api.RegisterView(iris.HTML(models.Config.App.Doc, ".html"))
	api.Handle("GET", "/v1/docs", func(ctx iris.Context) {
		ctx.View("index.html")
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

			admin.PartyFunc("/roles", func(roles router.Party) {
				roles.Get("/", controllers.GetAllRoles)
			})

			admin.PartyFunc("/perms", func(perms router.Party) {
				perms.Get("/", controllers.GetAllPerms)
			})

			admin.PartyFunc("/orders", func(orders router.Party) {
				orders.Get("/", controllers.GetAllOrders)
			})

			admin.PartyFunc("/clients", func(clients router.Party) {
				clients.Get("/", controllers.GetAllClients)
			})

			admin.PartyFunc("/plans", func(plans router.Party) {
				plans.Get("/", controllers.GetAllPlans)
				plans.Get("/parent", controllers.GetAllParentPlans)
			})

			admin.PartyFunc("/companies", func(companies router.Party) {
				companies.Get("/", controllers.GetAllCompanies)
			})

		})
	}
	api.Run(iris.Addr(":80"))
}
