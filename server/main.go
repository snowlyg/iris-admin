//go:generate go-bindata  ./www/dist/...
package main

import (
	"context"
	"fmt"
	"github.com/snowlyg/IrisAdminApi/server/seeder"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
	logger "github.com/sirupsen/logrus"
	"github.com/snowlyg/IrisAdminApi/server/config"
	"github.com/snowlyg/IrisAdminApi/server/libs"
	"github.com/snowlyg/IrisAdminApi/server/serve"
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
	irisServer *serve.Server
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

	opServer := serve.NewServer(Asset, AssetNames, AssetInfo)
	opServer.NewApp()
	prg.irisServer = opServer

	if len(os.Args) == 2 {
		if os.Args[1] == "version" {
			logger.Println(fmt.Sprintf("版本号：%s", Version))
			return
		} else if os.Args[1] == "seeder" {
			seeder.Run()
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
