package cron_server

import (
	"errors"
	"fmt"
	"sync"
	"time"

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

// DoOnce run job once time,this job will run after 2 second
func DoOnce(job cron.Job, t ...time.Duration) error {
	once := time.Now().Add(2 * time.Second)
	if len(t) == 1 {
		once = time.Now().Add(t[0] * time.Second)
	}
	onceSpec := fmt.Sprintf("%d %d %d %d %d %d", once.Second(), once.Minute(), once.Hour(), once.Day(), once.Month(), once.Weekday())
	_, err := CronInstance().AddJob(onceSpec, job)
	if err != nil {
		return err
	}
	return nil
}
