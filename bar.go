package admin

import (
	"log"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
)

type progressBar struct {
	bar *uiprogress.Bar
}

var steps = []string{
	"init gorm database",
	"init casbin authorization",
	"migrate database",
	"init locale",
	"sync router",
	"run server",
}

func newBar() *progressBar {
	bar := uiprogress.AddBar(len(steps)).AppendCompleted().PrependElapsed()
	pb := &progressBar{bar}
	go pb.deploy("init iris-admin")
	return pb
}

func (pb *progressBar) deploy(app string) {
	pb.bar.Width = 50

	// prepend the deploy step to the bar
	pb.bar.PrependFunc(func(b *uiprogress.Bar) string {
		if len(steps) > b.Current()-1 {
			info := app + ": " + steps[b.Current()-1]
			return strutil.Resize(info, uint(len(info)))
		}
		log.Printf("deploy current:%d %s\n", b.Current(), app)
		return ""
	})
}

func (pb *progressBar) Incr() {
	pb.bar.Incr()
	time.Sleep(time.Millisecond * 500)
}
