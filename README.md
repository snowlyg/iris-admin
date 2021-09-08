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

If you don't have a QQ account, you can into the [iris-go-tenancy/community](https://gitter.im/iris-go-tenancy/community?utm_source=share-link&utm_medium=link&utm_campaign=share-link) .

[![Gitter](https://badges.gitter.im/iris-go-tenancy/community.svg)](https://gitter.im/iris-go-tenancy/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) 
#### iris 学习记录分享

1.[Iris-go 项目登陆 API 构建细节实现过程](https://blog.snowlyg.com/iris-go-api-1/)

2.[iris + casbin 从陌生到学会使用的过程](https://blog.snowlyg.com/iris-go-api-2/)

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

#### 配合前端使用
- 编译前端页面到 admim 目录
```go
package main

import (
	"path/filepath"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/server/web"
)

func main() {
	webServer := web.Init()
	webServer.AddStatic("/", iris.Dir(filepath.Join(dir.GetCurrentAbPath(), "dist")))
	webServer.Run()
}
```

#### 简单用例
- [简单使用](https://github.com/snowlyg/IrisAdminApi/tree/master/example/simple)
- [添加模块](https://github.com/snowlyg/IrisAdminApi/tree/master/example/add_moule)
- [配合前端](https://github.com/snowlyg/IrisAdminApi/tree/master/example/single_with_vue)

#### 单元测试和接口文档[待更新]

#### 感谢 

[JetBrains](https://www.jetbrains.com/?from=IrisAdminApi) 对本项目的支持。

