<h1 align="center">IrisAdmin</h1>

<div align="center">
    <a href="https://travis-ci.org/snowlyg/IrisAdminApi"><img src="https://travis-ci.org/snowlyg/IrisAdminApi.svg?branch=blog" alt="Build Status"></a>
    <a href="https://codecov.io/gh/snowlyg/IrisAdminApi"><img src="https://codecov.io/gh/snowlyg/IrisAdminApi/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/report/github.com/snowlyg/IrisAdminApi"><img src="https://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/IrisAdminApi"><img src="https://godoc.org/github.com/snowlyg/IrisAdminApi?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/IrisAdminApi/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/IrisAdminApi" alt="Licenses"></a>
    <h5 align="center">IrisAdmin</h5>
</div>

> 简单项目仅供学习，欢迎指点！
>
#### 演示地址
主分支：
[http://irisadminapi.snowlyg.com](http://irisadminapi.snowlyg.com)

blog 分支:
[http://www.snowlyg.com](http://www.snowlyg.com) 

#### IRIS V12 中文文档
[IRIS V12 中文文档](https://www.snowlyg.com/chapter/1)


账号/密码 ： username/123456

###### `Iris-go` 学习交流 QQ 群 ：`676717248`
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=cc99ccf86be594e790eacc91193789746af7df4a88e84fe949e61e5c6d63537c"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Iris-go" title="Iris-go"></a>
---

#### 项目开发过程详解

1.[Iris-go 项目登陆 API 构建细节实现过程](https://www.snowlyg.com/#/detail/1)

2.[iris + casbin 从陌生到学会使用的过程](https://www.snowlyg.com/#/detail/2)

---

- 安装项目依赖

>加载依赖管理包 (解决国内下载依赖太慢问题)
>使用国内七牛云的 go module 镜像。
>
>参考 https://github.com/goproxy/goproxy.cn。
>
>阿里： https://mirrors.aliyun.com/goproxy/
>
>官方： https://goproxy.io/
>
>中国：https://goproxy.cn
>
>其他：https://gocenter.io

##### golang 1.13+ 可以直接执行：
```shell script
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

> 修改配置文件 `application.yml` ,配置文件需要放置在运行程序的同级目录

- 运行项目
>推荐使用 air 或者 gowatch 等热编译工具,直接使用 `go run main.go `  方法运行，可能会出现配置文件无法加载的问题
>如果想使用 `go run main.go` 命令运行

```shell script

# 安装工具 air     
go get -u github.com/cosmtrek/air

# 在 server 目录执行,可以通过 .air.conf 配置 air 工具
air
```

#### 登录项目
- http://localhost:80

#### 感谢 

[JetBrains](https://www.jetbrains.com/?from=IrisAdminApi) 对本项目的支持。

