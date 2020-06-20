<h1 align="center">IrisAdminApi</h1>

<div align="center">
    <a href="https://travis-ci.org/snowlyg/IrisAdminApi"><img src="https://travis-ci.org/snowlyg/IrisAdminApi.svg?branch=master" alt="Build Status"></a>
    <a href="https://codecov.io/gh/snowlyg/IrisAdminApi"><img src="https://codecov.io/gh/snowlyg/IrisAdminApi/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/report/github.com/snowlyg/IrisAdminApi"><img src="https://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/IrisAdminApi"><img src="https://godoc.org/github.com/snowlyg/IrisAdminApi?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/IrisAdminApi/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/IrisAdminApi" alt="Licenses"></a>
    <h5 align="center">IrisAdminApi</h5>
</div>

> 简单学习项目 ----写的挺烂，欢迎指点
>
#### 演示地址
[http://irisadminapi.snowlyg.com](http://irisadminapi.snowlyg.com)

账号/密码 ： username/123456


###### Iris-go 学习交流QQ群 ：676717248
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=cc99ccf86be594e790eacc91193789746af7df4a88e84fe949e61e5c6d63537c"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Iris-go" title="Iris-go"></a>

#### 新功能增加 ffmpeg 推流 (功能会继续更新)
- 通过 cgo 直接调用 ffmpeg C++ Api 的方式实现转码视频流,生成 .m3u8 文件到 hls 目录下。
- 配置文件增加了 recordpath: D:\Env\nginx\html\hls\cctv1 选项，需要配置绝对路径
- 启动项目后将地址 http://127.0.0.1:8085/record/out.m3u8 在 vlc 播放器打开即可播放中央9台。

#### 配置 ffmpeg api 库支持,增加系统变量
- windows 环境可以到 [github.com/snowlyg/ffmpegTest](github.com/snowlyg/ffmpegTest) 复制。
- 除了复制 lib,include 目录，并配置变量。还需要把 dll 目录下的所有 dll 文件复制到 backend 目录下。 
```shell script

export CGO_LDFLAGS="-L/usr/local/Cellar/ffmpeg/4.3_1/lib/ -lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter"
export CGO_CFLAGS="-I/usr/local/Cellar/ffmpeg/4.3_1/include"

```


![cctv9.png](cctv9.png)


#### 项目介绍
- `iris-go` 框架后台接口项目
- `gorm` 数据库模块 
- `jwt` 的单点登陆认证方式
- `cors` 跨域认证
- 数据支持 `mysql`，`sqlite3` 配置; `sqlite3` 需要下载 `gcc`, 并且在 `/temp` 目录下新建文件 `gorm.db` ,  `tgorm.db`。  [gcc 下载地址](http://mingw-w64.org/doku.php/download)
- 使用了 [https://github.com/snowlyg/gotransformer](https://github.com/snowlyg/gotransformer) 转换数据，返回数据格式化，excel 导入数据转换，xml 文件生产数据转换等 
- 增加了 `excel` 文件接口导入实例
- 前端采用了 `element-ui` 框架,代码集成到 `front` 目录
- 使用 `casbin` 做权限控制, `backend/config/rbac_model.conf` 为相关配置。系统会根据路由名称生成对应路由权限，并配置到管理员角色。
- 增加系统日志记录 `/logs` 文件夹下，自定义记录，控制器内 `ctx.Application().Logger().Infof("%s 登录系统",aul.Username)`

 **注意：**
 - 默认数据库设置为 `DriverType = "Sqlite"` ，使用 mysql 需要修改为 `DriverType = "Mysql"`，并且创建对应数据库 ,在 `backend/config/conf.tml` 文件中
 - `permissions.xlsx` 权限导入测试模板文件，仅供测试使用; 权限会自动生成，无需另外导入。
 
 -  `backend/config/config.go` 文件中的路径 `Root = os.Getenv("GOPATH") + "/src/github.com/snowlyg/IrisAdminApi/backend/"` 需要修改为你的项目路径,用于加载配置文件
 
---

#### 项目开发过程详解

[Iris-go 项目登陆 API 构建细节实现过程](https://learnku.com/articles/39551)
[iris + casbin 从陌生到学会使用的过程](https://learnku.com/articles/41416)

---


#### 更新日志

[UPDATE](UPDATE.MD)
---

#### 问题总结

[ERRORS](ERRORS.MD)



#### 项目初始化

>拉取项目

```shell script
git clone https://github.com/snowlyg/IrisAdminApi.git

# github 克隆太慢可以用 gitee 地址:

git clone https://gitee.com/snowlyg/IrisAdminApi.git

```

####  docker-compose 安装 （需要 docker 环境）

```shell script
  # 前端打包
    cd ./front
    npm install 
    npm run-script build

  # 复制配置文件，并修改配置
  # 复制到 config/ 目录即可。  docker-compose 脚本会将配置文件同步到 backend/config 目录下。
  cp config/application.yml.example config/application.yml 

  # 启动项目
  docker-compose up -d  
```



##### 普通环境安装项目

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
>
>golang 1.13 可以直接执行：
```shell script
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

```

>项目配置文件 backend/config/application.yml

```shell script
cp config/application.yml.example backend/config/application.yml
```

>打包前端代码 
```shell script
 cd front               # 进入前端代码目录
 npm install            #加载依赖
 npm run-script build   #打包前端代码

  # 复制前端文件到后端目录
  # 复制到 resources/app 到 backend/resources/app。
  cp -R resources/app backend/resources/app


 # 如果是开发前端代码,使用热加载
 npm run dev  
```


>运行项目 

[gowatch](https://gitee.com/silenceper/gowatch)
```shell script
go get github.com/silenceper/gowatch

# 安装 gowatch 后才可以使用
gowatch 

# 或者使用 go 命令（二选一）
go run main.go iris_base_rabc.go
```

---
##### 单元测试 
>http test

```shell script
 go test -v  //所有测试
 
 go test -run TestUserCreate -v //单个方法


// go get github.com/rakyll/gotest@latest 增加测试输出数据颜色

 gotest 
 
```

---

##### 接口文档
自动生成文档 (访问过接口就会自动成功)
因为原生的 jquery.min.js 里面的 cdn 是使用国外的，访问很慢。
有条件的可以开个 vpn ,如果没有可以根据下面的方法修改一下，访问就很快了
>打开 /resource/apiDoc/index.html 修改里面的

```shell script
https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js

国内的 cdn


https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```

>访问文档，从浏览器直接打开 http://localhost:8081/apiDoc

---

#### 登录项目
- http://localhost:8085

#### 感谢 

[JetBrains](https://www.jetbrains.com/?from=IrisAdminApi) 对本项目的支持。

