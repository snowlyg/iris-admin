package register

import (
	"IrisApiProject/controllers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

func Register(api *iris.Application) {
	api.Use(logger.New())

	api.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		_, _ = ctx.JSON(controllers.ApiResource(false, nil, "404 Not Found"))
	})
	api.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		_, _ = ctx.WriteString("Oups something went wrong, try again")
	})
}
