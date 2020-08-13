//go:generate go-bindata -prefix "assets" -fs  ./www/dist/...
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
	logger "github.com/sirupsen/logrus"
	"github.com/snowlyg/IrisAdminApi/server/config"
	"github.com/snowlyg/IrisAdminApi/server/libs"
	"github.com/snowlyg/IrisAdminApi/server/seeder"
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
	fmt.Println(fmt.Sprintf("host:%s", config.Config.Host))
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
	fmt.Println("start")
	go p.run()
	return nil
}

func (p *program) run() {
	fmt.Println("run")
	p.startIris()
}

func (p *program) stopIris() (err error) {
	if p.irisServer == nil {
		err = fmt.Errorf("HTTP Server Not Found")
		return
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = p.irisServer.App.Shutdown(ctx)
	if err != nil {
		return err
	}
	return
}

func (p *program) Stop(s service.Service) error {
	defer logger.Println("退出服务")
	defer f.Close()

	err := p.stopIris()
	if err != nil {
		fmt.Println(err.Error())
	}
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
		panic(err)
	}

	irisServer := serve.NewServer(AssetFile(), Asset, AssetNames)
	if irisServer == nil {
		panic("Http 初始化失败")
	}
	irisServer.NewApp()
	prg.irisServer = irisServer

	if len(os.Args) == 2 {
		if os.Args[1] == "version" {
			fmt.Println(fmt.Sprintf("版本号：%s", Version))
			return
		} else if os.Args[1] == "seeder" {
			seeder.Run()
			return
		} else if os.Args[1] == "perms" {
			fmt.Println("系统权限：")
			fmt.Println()
			routes := prg.irisServer.GetRoutes()
			for _, route := range routes {
				fmt.Println("+++++++++++++++")
				fmt.Println(fmt.Sprintf("名称 ：%s ", route.DisplayName))
				fmt.Println(fmt.Sprintf("路由地址 ：%s ", route.Name))
				fmt.Println(fmt.Sprintf("请求方式 ：%s", route.Act))
				fmt.Println()
			}

			return
		}

		err = service.Control(s, os.Args[1])
		if err != nil {
			panic(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		panic(err)
	}
}
