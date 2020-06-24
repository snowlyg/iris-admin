package main

import (
	stdContext "context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kardianos/service"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/backend/config"
	"github.com/snowlyg/IrisAdminApi/backend/libs"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	App    *iris.Application
	server *libs.Server
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())

	go func() {
		f := NewLogFile()
		defer f.Close()

		p.App = NewApp()
		p.App.Logger().SetOutput(f) //记录日志

		if config.Config.HTTPS {
			host := fmt.Sprintf("%s:%d", config.Config.Host, 443)
			if err := p.App.Run(iris.TLS(host, config.Config.Certpath, config.Config.Certkey)); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := p.App.Run(
				iris.Addr(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)),
				iris.WithoutServerError(iris.ErrServerClosed),
				iris.WithOptimizations,
				iris.WithTimeFormat(time.RFC3339),
			); err != nil {
				fmt.Println(err)
			}
		}
	}()

	go func() {
		fmt.Println("hls start")
		p.server = libs.GetServer()
		_ = p.server.Start()
	}()

	return nil
}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")

	p.shutdown()
	p.server.Stop()

	return nil
}

func (p *program) shutdown() {
	time.Sleep(3 * time.Second)
	ctx, cancel := stdContext.WithTimeout(stdContext.TODO(), 3*time.Second)
	defer cancel()
	p.App.Shutdown(ctx)
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {

	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "backend",
		DisplayName: "Go Service Example for Logging",
		Description: "This is an example Go service that outputs log messages.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
