package web

import (
	"testing"

	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
)

// WebFunc 框架服务接口
// - GetTestClient 测试客户端
// - GetTestLogin 测试登录
// - AddWebStatic 添加静态页面
// - InitDriver 初始化认证
// - AddUploadStatic 上传文件路径
// - Run 启动
type WebFunc interface {
	GetTestClient(t *testing.T) *tests.Client
	GetTestLogin(t *testing.T, url string, res tests.Responses, datas ...interface{}) *tests.Client

	AddWebStatic()
	AddUploadStatic()
	InitDriver() error
	InitRouter() error
	Run()
}

// Start 启动 web 服务
func Start(wf WebFunc) {
	err := wf.InitDriver()
	if err != nil {
		zap_server.ZAPLOG.Error("初始化系统失败", zap.String("wf.InitDriver", err.Error()))
		return
	}
	err = wf.InitRouter()
	if err != nil {
		zap_server.ZAPLOG.Error("初始化路由失败", zap.String("wf.InitRouter", err.Error()))
		return
	}
	wf.Run()
}

// StartTest 启动 web 服务
func StartTest(wf WebFunc) {
	err := wf.InitDriver()
	if err != nil {
		zap_server.ZAPLOG.Error("初始化系统失败", zap.String("wf.InitDriver", err.Error()))
		return
	}
	err = wf.InitRouter()
	if err != nil {
		zap_server.ZAPLOG.Error("初始化路由失败", zap.String("wf.InitRouter", err.Error()))
		return
	}
}
