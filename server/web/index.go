package web

import (
	"fmt"
	"log"

	"github.com/snowlyg/iris-admin/server/viper_server"
)

// init
func init() {
	viper_server.Init(getViperConfig())
}

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
func Start(wf WebFunc) error {
	if err := wf.InitRouter(); err != nil {
		return fmt.Errorf("init router fail:%s", err.Error())
	}
	wf.Run()
	return nil
}

// StartTest
func StartTest(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		log.Printf("start test fail:%s\n", err.Error())
	}
}
