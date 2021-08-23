package web

import (
	"github.com/kataras/iris/v12"
)

type WebModule struct {
	relativePath string
	handler      func(p iris.Party)
}

func NewModule(relativePath string, handler func(index iris.Party)) WebModule {
	return WebModule{
		relativePath: relativePath,
		handler:      handler,
	}
}
