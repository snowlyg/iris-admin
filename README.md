<h1 align="center">IrisAdminApi</h1>

<p align="center">
    <a href="https://travis-ci.org/snowlyg/IrisAdminApi"><img src="https://travis-ci.org/snowlyg/IrisAdminApi.svg?branch=master" alt="Build Status"></a>
    <a href="https://codecov.io/gh/snowlyg/IrisAdminApi"><img src="https://codecov.io/gh/snowlyg/IrisAdminApi/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/report/github.com/snowlyg/IrisAdminApi"><img src="https://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/IrisAdminApi"><img src="https://godoc.org/github.com/snowlyg/IrisAdminApi?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/IrisAdminApi/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/IrisAdminApi" alt="Licenses"></a>
    <h5 align="center">Iris后台接口项目</h5>
</p>

#### 项目介绍
- iris 框架后台接口项目
- 采用了 gorm 数据库模块 和 jwt 的单点登陆认证方式
- 测试默认使用了 sqlite3 数据库
---

#### 项目更新
- 修改项目配置加载方式
- 简化 tools、database 使用
- 增加数据格式化返回 [snowlyg/gotransformer](https://github.com/snowlyg/gotransformer)
---

#### api项目初始化

>拉取项目

```
git clone https://github.com/snowlyg/IrisAdminApi.git
```
//github 太慢可以用 gitee
```
git clone https://gitee.com/dtouyu/IrisAdminApi.git
```

>加载依赖管理包 (解决国内下载依赖太慢问题)
>使用国内七牛云的 go module 镜像。
 参考 https://github.com/goproxy/goproxy.cn。
>golang 1.13 可以直接执行：
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

// 注意 goland 会覆盖 GOPROXY 配置，建议更新依赖包的时候关闭 `goland` 的`module`配置，更新完再打开。
```


>项目配置文件 /config/config.toml

```
cp config.toml.example config.toml
```

>运行项目 

[gowatch](https://gitee.com/silenceper/gowatch)
```
go get github.com/silenceper/gowatch

gowatch //安装 gowatch 后才可以使用

go run main.go // go 命令
```

cd github.com/silenceper/gowatch
go install 

```



---
##### 单元测试 
>http test

```
 go test -v  //所有测试
 
 go test -run TestUserCreate -v //单个测试
 
```

---

##### api 文档使用
自动生成文档 (访问过接口就会自动成功)
因为原生的 jquery.min.js 里面的 cdn 是使用国外的，访问很慢。
有条件的可以开个 vpn ,如果没有可以根据下面的方法修改一下，访问就很快了
>打开 /resource/apiDoc/index.html 修改里面的

```
https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js

国内的 cdn


https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```

>访问文档，从浏览器直接打开 http://localhost/apiDoc

---

#### 登录项目
输入地址 http://localhost:8081

//在 conig/conf.tml 内配置 

项目管理员账号 ： username
项目管理员密码 ： password

##### 问题总结

[问题总结](ERRORS.MD)
