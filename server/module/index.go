package module

import (
	"github.com/kataras/iris/v12"
)

type InitDBFunc interface {
	Init() (err error)
}

type WebModule struct {
	RelativePath string
	Handler      func(p iris.Party)
	Modules      []WebModule
}

func NewModule(relativePath string, handler func(index iris.Party), modules ...WebModule) WebModule {
	return WebModule{
		RelativePath: relativePath,
		Handler:      handler,
		Modules:      modules,
	}
}

func (wm *WebModule) GetModules() []WebModule {
	return wm.Modules
}
