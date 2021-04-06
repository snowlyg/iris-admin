<h1 align="center">IrisAdmin</h1>

<div align="center">
    <a href="https://codecov.io/gh/snowlyg/IrisAdminApi"><img src="https://codecov.io/gh/snowlyg/IrisAdminApi/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/report/github.com/snowlyg/IrisAdminApi"><img src="https://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/IrisAdminApi"><img src="https://godoc.org/github.com/snowlyg/IrisAdminApi?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/IrisAdminApi/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/IrisAdminApi" alt="Licenses"></a>
</div>

> 简单项目仅供学习，欢迎指点！

[IRIS V12 中文文档](https://github.com/snowlyg/iris/wiki)


###### `Iris-go` 学习交流 QQ 群 ：`676717248`
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=cc99ccf86be594e790eacc91193789746af7df4a88e84fe949e61e5c6d63537c"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Iris-go" title="Iris-go"></a>
---

#### 项目开发过程详解

1.[Iris-go 项目登陆 API 构建细节实现过程](https://blog.snowlyg.com/iris-go-api-1/)

2.[iris + casbin 从陌生到学会使用的过程](https://blog.snowlyg.com/iris-go-api-2/)

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
- 复制配置文件
```
cp application.example.yml application.yml
```

>  修改配置文件 `application.yml` 

- 运行项目
>如果想使用 `go run main.go --config ` 命令运行,注意不用 --config 指定配置路径，将无法加载配置文件
```
# your_config_path 指定配置文件绝对路径
 go run main.go --config your_config_path
```

>推荐使用 air 热编译工具
```
# 安装工具 air     
go get -u github.com/cosmtrek/air

cp .air.example.conf  .air.conf # 复制后修改 .air.conf 文件，默认为 mac 环境

air
```

- 填充数据, 注意配置文件同项目配置文件，权限数据位于 tools/seed/data
```
go build -o seed tools/seed/main.go 
# youer_seed_data_path 指定目录即可
./seed --config your_config_path --path youer_seed_data_path
```

#### 报错 Error 1071: Specified key was too long; max key length is 1000 bytes
- 修改数据库引擎为 InnoDB

#### postman 接口
```text
https://www.getpostman.com/collections/048078cdfd16667352b0
```

#### 运行测试
```
go test ./...
```

#### 感谢 

[JetBrains](https://www.jetbrains.com/?from=IrisAdminApi) 对本项目的支持。

