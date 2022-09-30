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

// WebFunc
// - GetTestClient
// - GetTestLogin
// - AddWebStatic
// - AddUploadStatic
// - Run
type WebFunc interface {
	WebBaseFunc
}

// Start
func Start(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return
	}
	wf.Run()
}

// StartTest
func StartTest(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
	}
}

// InitWeb
func InitWeb() error {
	err := viper_server.Init(getViperConfig())
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return err
	}
	return nil
}
