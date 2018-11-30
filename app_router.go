package main

import (
	"IrisApiProject/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
	"net/http"
	"time"
)

/**
 * 判断 token 是否有效
 * 如果有效 就获取信息并且保存到请求里面
 * @method AuthToken
 * @param  {[type]}  ctx       iris.Context    [description]
 */
func AuthToken(ctx iris.Context) {
	u := ctx.Values().Get("jwt").(*jwt.Token) //获取 token 信息
	token := MGetOauthTokenByToken(u.Raw)     //获取 access_token 信息
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(ApiJson{Status: false, Data: "", Msg: "token 已经过期"})
		ctx.Next()

		return
	}

	user := new(Users)
	user.ID = token.UserId

	user.GetUserById() //获取 user 信息

	ctx.Values().Set("auth_user_id", user.ID)
	ctx.Values().Set("auth_user_name", user.Name)

	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}

/**
 * 初始化 iris app
 * @method NewApp
 * @return  {[type]}  api      *iris.Application  [iris app]
 */
func NewApp() (api *iris.Application) {

	api = iris.New()
	api.Use(logger.New())

	api.OnErrorCode(iris.StatusNotFound, NotFound)
	api.OnErrorCode(iris.StatusInternalServerError, InternalServerError)

	db.AutoMigrate(new(Users), new(Settings), new(OauthToken), new(Orders), new(Companies), new(Roles))

	iris.RegisterOnInterrupt(func() {
		db.Close()
	})

	//"github.com/iris-contrib/middleware/cors"
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	appName := conf.Get("app.name").(string)
	appDoc := conf.Get("app.doc").(string)
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: appName,
		DocPath:  appDoc + "/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": conf.Get("app.url").(string),
			"Staging":    "",
		},
	})

	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Use(middleware.NewYaag()) // <- IMPORTANT, register the middleware.
		v1.Post("/admin/login", CUserLogin)
		v1.PartyFunc("/admin", func(admin router.Party) {
			admin.Use(middleware.JwtHandler().Serve, AuthToken)
			admin.Get("/", CGetHomeData)
			admin.Get("/logout", CUserLogout)

			admin.PartyFunc("/users", func(users router.Party) {
				users.Get("/", CGetAllUsers)
				users.Get("/{id:uint}", CGetUser)
				users.Post("/", CCreateUser)
				users.Post("/{id:uint}/update", CUpdateUser)
				users.Get("/{id:uint}/frozen", CFrozenUser)
				users.Get("/{id:uint}/audit", CSetUserAudit)
				users.Get("/{id:uint}/refrozen", CRefrozenUser)
				users.Delete("/{id:uint}", CDeleteUser)
				users.Get("/profile", CGetProfile)
			})

			admin.PartyFunc("/roles", func(roles router.Party) {
				roles.Get("/", CGetAllRoles)
			})

			admin.PartyFunc("/perms", func(perms router.Party) {
				perms.Get("/", CGetAllPerms)
			})

			admin.PartyFunc("/orders", func(orders router.Party) {
				orders.Get("/", CGetAllOrders)
			})

			admin.PartyFunc("/clients", func(clients router.Party) {
				clients.Get("/", CGetAllClients)
			})

			admin.PartyFunc("/plans", func(plans router.Party) {
				plans.Get("/", CGetAllPlans)
				plans.Get("/parent", CGetAllParentPlans)
			})

			admin.PartyFunc("/companies", func(companies router.Party) {
				companies.Get("/", CGetAllCompanies)
			})

		})
	}

	return
}
