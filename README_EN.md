<h1 align="center">IrisAdmin</h1>

<div align="center">
    <a href="https://codecov.io/gh/snowlyg/IrisAdminApi"><img src="https://codecov.io/gh/snowlyg/IrisAdminApi/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="hhttps://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi"><img src="https://goreportcard.com/badge/github.com/snowlyg/IrisAdminApi" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/IrisAdminApi"><img src="https://godoc.org/github.com/snowlyg/IrisAdminApi?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/iris-admin/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/IrisAdminApi" alt="Licenses"></a>
</div>

[简体中文](./README.md)  | English

#### Project url
[GITHUB](https://github.com/snowlyg/iris-admin) | [GITEE](https://gitee.com/snowlyg/iris-admin) 
****
> This project just for learning golang, welcome to give your suggestions!

#### Documentation
- [IRIS V12 document for chinese](https://github.com/snowlyg/iris/wiki)
- [godoc](https://pkg.go.dev/github.com/snowlyg/iris-admin?utm_source=godoc)

#### COMMUNICATIONS
- `Iris-go`  QQ group for communication：`676717248`
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=cc99ccf86be594e790eacc91193789746af7df4a88e84fe949e61e5c6d63537c"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Iris-go" title="Iris-go"></a>

- If you don't have a QQ account, you can into the [iris-go-tenancy/community](https://gitter.im/iris-go-tenancy/community?utm_source=share-link&utm_medium=link&utm_campaign=share-link) [![Gitter](https://badges.gitter.im/iris-go-tenancy/community.svg)](https://gitter.im/iris-go-tenancy/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) .
- WeChat group please add WeChat ID: `c25vd2x5Z19jaGluYQ==`


#### BLOG

- [REST API with iris-go web framework ](https://blog.snowlyg.com/iris-go-api-1/)

- [How to user iris-go with casbin](https://blog.snowlyg.com/iris-go-api-2/)

---

#### Getting started

```go
package main

import (
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
)

func main() {
	web.Start(web_iris.Init())
}
```

#### Run project 
- When you first run this cmd `go run main.go` , you can see some config files in  the `config` directory,
- and `rbac_model.conf` will be created in your project root directory.
```sh
go run main.go
```

#### Module
- The framework has a built-in v1 version of the basic authentication module by default.
- Your can use AddModule() to add other modules .
```go
package main

import (
	v1 "github.com/snowlyg/iris-admin/modules/v1"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
)

func main() {
	wi := web_iris.Init()
	v1Party := web_iris.Party{
		Perfix:    "/api/v1",
		PartyFunc: v1.Party(),
	}
	wi.AddModule(v1Party)
	web.Start(web_iris.Init())
}
```

#### Default static file path
- A static file access path has been built in by default
- Static files will upload to `/static/upload` directory.
- You can set this config key `static-path` to change the default directory.
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

#### Add Static file path
- You can add static file access path,through `AddStatic` function.
```go
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

func main() {
	webServer := web_iris.Init()
    fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), "/other"))
	webServer.AddStatic("/other",fsOrDir)
	webServer.Run()
}
```

#### Use with front-end framework , e.g. vue.
- Default,you must build vue to the `dist` directory.
- Naturally you can set this config key `web-path` to change the default directory.
```go
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

func main() {
	webServer := web_iris.Init()
	webServer.AddWebStatic("/")
	webServer.Run()
}
```
- Front-end page reference/borrowing:
 *notice: The front-end only realizes preview effect simply*
- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin/tree/master/web)
- [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)


#### Example
- [example](https://github.com/snowlyg/iris-admin/tree/master/example)

#### Unit test and documentation: [to be updated]
- Before start unit tests,you must create two files which is named `redis_pwd.txt `and `redis_pwd.txt` ,fill `redis` and `mysql` 's password in to these two files separately, on the directory the `main_test.go` file is located.
- [helper/tests](https://github.com/snowlyg/helper/tree/main/tests) package the unit test used, it's  simple package base on [httpexpect/v2](https://github.com/gavv/httpexpect).
- [example for unit test](https://github.com/snowlyg/iris-admin/tree/master/modules/v1/user/test)

#### Thanks 

 - Thanks [JetBrains](https://www.jetbrains.com/?from=IrisAdminApi)' supports .


