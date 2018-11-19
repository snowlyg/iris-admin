package controllers

import "github.com/kataras/iris"

func NotFound(ctx iris.Context)  {
	ctx.JSON(iris.Map{"status":false,"msg":"404 not found"})
}
func InternalServerError(ctx iris.Context) {
	ctx.WriteString("Oups something went wrong, try again")
}