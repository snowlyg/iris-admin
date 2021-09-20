package module

import (
	"github.com/kataras/iris/v12"
)

// InitDBFunc 数据化初始化接口
type InitDBFunc interface {
	Init() (err error)
}

// WebModule web 模块结构
// - RelativePath 关联路径
// - Handler 模块 Handler
// - Modules 子模块
type WebModule struct {
	RelativePath string
	Handler      func(p iris.Party)
	Modules      []WebModule
}

// NewModule 添加新的模块
func NewModule(relativePath string, handler func(index iris.Party), modules ...WebModule) WebModule {
	return WebModule{
		RelativePath: relativePath,
		Handler:      handler,
		Modules:      modules,
	}
}

// GetModules 获取模块列表
func (wm *WebModule) GetModules() []WebModule {
	return wm.Modules
}
