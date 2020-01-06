package register

import (
	"IrisApiProject/controllers"
	"github.com/kataras/iris/v12/middleware/logger"

	"github.com/kataras/iris/v12"
	//"github.com/casbin/casbin/v2"
	//cm "github.com/iris-contrib/middleware/casbin"
)

//var Enforcer, _ = casbin.NewEnforcer("casbinmodel.conf", "casbinpolicy.csv")
func Register(api *iris.Application) {
	api.Use(logger.New())
	//casbinMiddleware := cm.New(Enforcer)
	//api.Use(casbinMiddleware.ServeHTTP)
	api.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		_, _ = ctx.JSON(controllers.ApiResource(false, nil, "404 Not Found"))
	})
	api.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		_, _ = ctx.WriteString("Oups something went wrong, try again")
	})
}
