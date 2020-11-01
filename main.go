//go:generate go-bindata -prefix "./views" -fs  ./views/...
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
	logger "github.com/sirupsen/logrus"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/seeder"
	"github.com/snowlyg/blog/web_server"
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
	host := fmt.Sprintf("%s:%d", libs.Config.Host, libs.Config.Port)
	if host != "" {
		go func() {
			logger.Println("HTTP-IRIS listen On ", host)
			err := p.irisServer.Serve()
			if err != nil {
				panic(err)
			}
		}()
	}
}

type program struct {
	irisServer *web_server.Server
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
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
		Name:        "GoIrisAdminApi",
		DisplayName: "GoIrisAdminApi",
		Description: "go+web+iris",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}

	irisServer := web_server.NewServer(AssetFile())
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
		} else if os.Args[1] == "sync_perms" {
			seeder.AddPerm()
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
		} else if os.Args[1] == "start" {
			if libs.IsPortInUse(libs.Config.Port) {
				panic(fmt.Sprintf("端口 %d 已被使用", libs.Config.Port))
			}
		}

		err = service.Control(s, os.Args[1])
		if err != nil {
			panic(err)
		}
		return
	}

	if libs.IsPortInUse(libs.Config.Port) {
		panic(fmt.Sprintf("端口 %d 已被使用", libs.Config.Port))
	}

	err = s.Run()
	if err != nil {
		panic(err)
	}
}
