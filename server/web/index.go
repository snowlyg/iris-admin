package web

import (
	"testing"

	"github.com/snowlyg/helper/tests"
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
	GetTestLogin(t *testing.T, url string, res tests.Responses, datas ...map[string]interface{}) *tests.Client
	AddWebStatic(perfix string)
	AddUploadStatic()
	InitDriver() error
	InitRouter() error
	Run()
}

// Start 启动 web 服务
func Start(wf WebFunc) {
	err := wf.InitDriver()
	if err != nil {
		panic(err)
	}
	err = wf.InitRouter()
	if err != nil {
		panic(err)
	}
	wf.Run()
}
