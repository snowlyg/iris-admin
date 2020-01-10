package register

import (
	"net/http"

	"github.com/kataras/iris/v12"
	//"github.com/casbin/casbin/v2"
	//cm "github.com/iris-contrib/middleware/casbin"
)

func Register(api *iris.Application) {
	api.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.StatusCode(http.StatusNotFound)
		ctx.Next()
		return
	})
	api.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		_, _ = ctx.WriteString("Oups something went wrong, try again")
	})
}
