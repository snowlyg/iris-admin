package admin

import (
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
)

type progressBar struct {
	// sync.WaitGroup
	bar *uiprogress.Bar
}

var steps = []string{
	"init gorm database",
	"init casbin authorization",
	"migrate database",
	"init locale",
	"run server",
}

func newBar() *progressBar {
	bar := uiprogress.AddBar(len(steps)).AppendCompleted().PrependElapsed()
	pb := &progressBar{bar}
	go pb.deploy("init iris-admin")
	return pb
}

func (pb *progressBar) deploy(app string) {
	// defer pb.Done()
	pb.bar.Width = 50

	// prepend the deploy step to the bar
	pb.bar.PrependFunc(func(b *uiprogress.Bar) string {
		info := app + ": " + steps[b.Current()-1]
		return strutil.Resize(info, uint(len(info)))
	})

	// rand.New(rand.NewSource(500))
	// for pb.bar.Incr() {
	// 	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	// }
}

func (pb *progressBar) Incr() {
	pb.bar.Incr()
	time.Sleep(time.Millisecond * 500)
}
