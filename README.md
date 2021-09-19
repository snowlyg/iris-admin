<h1 align="center">IrisAdmin</h1>

<div align="center">
    <a href="https://codecov.io/gh/snowlyg/IrisAdminApi"><img src="https://codecov.io/gh/snowlyg/IrisAdminApi/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/report/github.com/snowlyg/IrisAdminApi"><img src="https://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/IrisAdminApi"><img src="https://godoc.org/github.com/snowlyg/IrisAdminApi?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/iris-admin/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/IrisAdminApi" alt="Licenses"></a>
</div>

简体中文 | [English](./README_EN.md) 

#### 项目地址
[GITHUB](https://github.com/snowlyg/iris-admin) | [GITEE](https://gitee.com/snowlyg/iris-admin) 

> 简单项目仅供学习，欢迎指点！

#### 相关文档
- [IRIS V12 中文文档](https://github.com/snowlyg/iris/wiki)
- []()

#### 交流方式
- `Iris-go` 学习交流 QQ 群 ：`676717248`
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=cc99ccf86be594e790eacc91193789746af7df4a88e84fe949e61e5c6d63537c"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Iris-go" title="Iris-go"></a>

- If you don't have a QQ account, you can into the [iris-go-tenancy/community](https://gitter.im/iris-go-tenancy/community?utm_source=share-link&utm_medium=link&utm_campaign=share-link) .
- 微信群请加微信号: `c25vd2x5Z19jaGluYQ==`

[![Gitter](https://badges.gitter.im/iris-go-tenancy/community.svg)](https://gitter.im/iris-go-tenancy/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) 
#### iris 学习记录分享

- [Iris-go 项目登陆 API 构建细节实现过程](https://blog.snowlyg.com/iris-go-api-1/)

- [iris + casbin 从陌生到学会使用的过程](https://blog.snowlyg.com/iris-go-api-2/)

---

#### 简单使用
```go
package main

import (
	"github.com/snowlyg/iris-admin/server/web"
)

func main() {
	webServer := web.Init()
	webServer.Run()
}
```

#### 启动项目
- 第一次启动项目后,会自动生成 `config.yaml` 和 `rbac_model.conf` 两个配置文件
```sh
go run main.go
```

#### 添加模块
- 框架默认内置了v1 版本的基础认证模块
- 可以使用 AddModule() 增加其他 admin模块
```go
package main

import (
  	"github.com/snowlyg/iris-admin/server/web"
  "github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party admin模块
func Party() module.WebModule {
  handler := func(admin iris.Party) {
    // 中间件
    admin.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		admin.Get("/", GetAllAdmins).Name = "admin列表"
	}
	return module.NewModule("/admins", handler)
}

func GetAllAdmins(ctx iris.Context) {
  // 处理业务逻辑
  // ... 
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: list, Msg: g.NoErr.Msg})
}

func main() {
	webServer := web.Init()
    webServer.AddModule(Party())
	webServer.Run()
}
```

#### 设置静态文件路径
- 已经默认内置了一个静态文件访问路径
- 静态文件将会上传到 `/static/upload` 目录
- 可以修改配置项 `static-path` 修改默认目录
```yaml
system:
  addr: 127.0.0.1:8085
  cache-type: ""
  db-type: ""
  level: debug
  static-path: /static/upload
  static-prefix: /upload
  time-format: "2006-01-02 15:04:05"
  web-path: ./dist
```

#### 设置其他静态文件路径
- 设置其他静态文件路径，可以使用 `AddStatic` 方法
```go
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

func main() {
	webServer := web.Init()
    fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), "/other"))
	webServer.AddStatic("/other",fsOrDir)
	webServer.Run()
}
```

#### 配合前端使用
- 编译前端页面默认 `dist` 目录
- 可以修改配置项 `web-path` 修改默认目录
```go
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

func main() {
	webServer := web.Init()
	webServer.AddWebStatic("/")
	webServer.Run()
}
```
- 前端页面参考/借用：【前端只简单实现预览效果】
- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin/tree/master/web)
- [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)


#### 简单用例
- [简单使用](https://github.com/snowlyg/iris-admin/tree/master/example)

#### 单元测试和接口文档[待更新] 
- 测试前在 `main_test.go` 文件所在目录新建 `redis_pwd.txt `和 `redis_pwd.txt` 两个文件,分别填入 `redis` 和 `mysql` 的密码
- 测试使用依赖库 [helper/tests](https://github.com/snowlyg/helper/tree/main/tests) 是基于 [httpexpect/v2](https://github.com/gavv/httpexpect) 的简单封装
- [接口单元测试例子](https://github.com/snowlyg/iris-admin/tree/master/modules/v1/user/test)

#### 感谢 

[JetBrains](https://www.jetbrains.com/?from=IrisAdminApi) 对本项目的支持。


