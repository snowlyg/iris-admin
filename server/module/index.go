package module

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/str"
)

// InitSourceFunc 数据化初始化接口
type InitSourceFunc interface {
	Init() (err error)
}

// WebModule web 模块结构
// - RelativePath 关联路径
// - Handler 模块 Handler
// - Modules 子模块
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

// GetModules 获取子模块列表
func (wm *WebModule) GetModules() []WebModule {
	return wm.Modules
}

// GetMigrateId 获取迁移id
func GetMigrateId(perfix string) string {
	return str.Join(time.Now().Format("2006_01_02_15_04_05_"), perfix)
}
