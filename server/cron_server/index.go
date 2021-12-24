package cronserver

import (
	"errors"
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	ErrNotLocalServer = errors.New("不能控制非本机服务")
	ErrStartServer    = errors.New("服务已经启动")
)

var (
	once sync.Once
	cc   *cron.Cron
)

// CronInstance cron 单例
func CronInstance() *cron.Cron {
	once.Do(func() {
		cc = cron.New(cron.WithSeconds())
	})
	return cc
}
