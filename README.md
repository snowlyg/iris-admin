<h1 align="center">IrisAdminApi</h1>

<div align="center">
    <a href="https://travis-ci.org/snowlyg/IrisAdminApi"><img src="https://travis-ci.org/snowlyg/IrisAdminApi.svg?branch=master" alt="Build Status"></a>
    <a href="https://codecov.io/gh/snowlyg/IrisAdminApi"><img src="https://codecov.io/gh/snowlyg/IrisAdminApi/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/report/github.com/snowlyg/IrisAdminApi"><img src="https://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/IrisAdminApi"><img src="https://godoc.org/github.com/snowlyg/IrisAdminApi?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/IrisAdminApi/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/IrisAdminApi" alt="Licenses"></a>
    <h5 align="center">Iris后台接口项目</h5>
</div>

#### 项目介绍
- iris golang框架后台接口项目
- gorm 数据库模块 
- jwt 的单点登陆认证方式
- cors 跨域认证
- 测试使用了 sqlite3 数据库 (需要 gcc 环境，如果不需要单元测试可以直接注释 sqlite3 依赖，或者采用 mysql 测试)
- 使用了 [https://github.com/snowlyg/gotransformer](https://github.com/snowlyg/gotransformer) 转换数据，返回数据格式化，excel 导入数据转换，xml 文件生产数据转换等 
- 增加了 excel 文件接口导入实例 (数据导入效率比 PHP 高很多，本人只尝试了 php 和 go 之间的对比)
- 前端采用了 element-ui 框架，如果需要修改前端代码请移步 [https://github.com/snowlyg/IrisApiVueAdmin](https://github.com/snowlyg/IrisApiVueAdmin)
- 使用 casbin 做权限控制, config/rbac_model.conf 为相关配置。系统会根据路由名称生成对应路由权限，并配置到管理员角色。
- **注意：更新代码后，需要手动删所有数据表，数据结构有调整。**

---

#### 更新日志
[更新日志](UPDATE.MD)
---

#### 问题总结
[问题记录](ERRORS.MD)


#### 项目初始化

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
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

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

---
##### 单元测试 
>http test

```
 go test -v  //所有测试
 
 go test -run TestUserCreate -v //单个方法
 
```

---

##### 接口文档
自动生成文档 (访问过接口就会自动成功)
因为原生的 jquery.min.js 里面的 cdn 是使用国外的，访问很慢。
有条件的可以开个 vpn ,如果没有可以根据下面的方法修改一下，访问就很快了
>打开 /resource/apiDoc/index.html 修改里面的

```
https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js

国内的 cdn


https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```

>访问文档，从浏览器直接打开 http://localhost:8081/apiDoc

---

#### 登录项目
输入地址 http://localhost:8081

//在 conig/conf.tml 内配置 

项目管理员账号 ： username
项目管理员密码 ： password



###### Iris-go 学习交流QQ群 ：676717248
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=cc99ccf86be594e790eacc91193789746af7df4a88e84fe949e61e5c6d63537c"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Iris-go" title="Iris-go"></a>

