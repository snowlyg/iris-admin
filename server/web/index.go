package web

import (
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

type WebBaseFunc interface {
	AddWebStatic(staticAbsPath, webPrefix string, paths ...string)
	AddUploadStatic(staticAbsPath, webPrefix string)
	InitRouter() error
	Run()
}

// WebFunc 框架服务接口
// - GetTestClient 测试客户端
// - GetTestLogin 测试登录
// - AddWebStatic 添加静态页面
// - AddUploadStatic 上传文件路径
// - Run 启动
type WebFunc interface {
	WebBaseFunc
}

// Start 启动服务
func Start(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return
	}
	wf.Run()
}

// StartTest 启动服务
func StartTest(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
	}
}

// InitWeb 初始化配置
func InitWeb() error {
	err := viper_server.Init(getViperConfig())
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return err
	}
	return nil
}
