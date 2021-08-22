package web

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type WebModule struct {
	relativePath string
	handler      func(p iris.Party)
	middlewares  []context.Handler //中间件
}

func NewModule(relativePath string, handler func(index iris.Party), middlewares ...context.Handler) WebModule {
	return WebModule{
		relativePath: relativePath,
		handler:      handler,
		middlewares:  middlewares,
	}
}
