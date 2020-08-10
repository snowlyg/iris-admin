package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
	logger "github.com/sirupsen/logrus"
	"github.com/snowlyg/IrisAdminApi/backend/config"
	"github.com/snowlyg/IrisAdminApi/backend/libs"
	"github.com/snowlyg/IrisAdminApi/backend/server"
)

var Version = "master"
var f *os.File

func init() {
	f, _ := os.OpenFile(fmt.Sprintf("%s/go_iris.log", libs.LogDir()), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	logger.SetFormatter(&logger.JSONFormatter{})
	logger.SetOutput(f)
	logger.SetLevel(logger.DebugLevel)
}

func (p *program) startIris() {
	host := config.Config.Host
	if host != "" {
		go func() {
			logger.Println("HTTP-IRIS listen On ", host)
			err := p.irisServer.Serve()
			if err != nil {
				logger.Println("HTTP-IRIS listen Err :", err)
			}
		}()
	}
}

type program struct {
	irisServer *server.Server
}

func (p *program) Start(s service.Service) error {
	go p.run()

	p.startIris()
	return nil
}
func (p *program) run() {

}

func (p *program) stopIris() (err error) {
	if p.irisServer == nil {
		err = fmt.Errorf("HTTP Server Not Found")
		return
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p.irisServer.App.Shutdown(ctx)
	return
}

func (p *program) Stop(s service.Service) error {
	defer logger.Println("退出服务")
	defer f.Close()

	p.stopIris()
	return nil
}

func main() {
	log.Println(fmt.Sprintf(` 
============================================
   ___  ___ ___ ___ ___ ___   _   ___ ___ 
  / __|/ _ \_ _| _ \_ _/ __| /_\ | _ \_ _|
 | (_ | (_) | ||   /| |\__ \/ _ \|  _/| | 
  \___|\___/___|_|_\___|___/_/ \_\_| |___|

============================================

version: %s`, Version))

	svcConfig := &service.Config{
		Name:        "GoIrisApi",
		DisplayName: "GoIrisApi",
		Description: "go+web+iris 后台服务",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Println(err)
	}

	opServer := server.NewServer()
	opServer.NewApp()
	prg.irisServer = opServer

	if len(os.Args) == 2 {
		if os.Args[1] == "version" {
			logger.Println(fmt.Sprintf("版本号：%s", Version))
			return
		}

		err = service.Control(s, os.Args[1])
		if err != nil {
			logger.Fatal(err)
		}
		return

	}
	err = s.Run()
	if err != nil {
		logger.Println(err)
	}
}
