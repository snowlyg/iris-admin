/*
<h1 align="center">IrisAdmin</h1>

<div align="center">
    <a href="https://codecov.io/gh/snowlyg/iris-admin"><img src="https://codecov.io/gh/snowlyg/iris-admin/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/badge/github.com/snowlyg/iris-admin"><img src="https://goreportcard.com/badge/github.com/snowlyg/iris-admin" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/iris-admin"><img src="https://godoc.org/github.com/snowlyg/iris-admin?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/iris-admin/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/iris-admin" alt="Licenses"></a>
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
-  [iris-go-tenancy/community](https://gitter.im/iris-go-tenancy/community?utm_source=share-link&utm_medium=link&utm_campaign=share-link) [![Gitter](https://badges.gitter.im/iris-go-tenancy/community.svg)](https://gitter.im/iris-go-tenancy/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) .


#### BLOG

- [REST API with iris-go web framework ](https://blog.snowlyg.com/iris-go-api-1/)

- [How to user iris-go with casbin](https://blog.snowlyg.com/iris-go-api-2/)

---


#### Program introduction

##### The project consists of multiple services, each with different functions.

- [viper_server]
- - The service configuration is initialized and generate a local configuration file.
- - Use [github.com/spf13/viper](https://github.com/spf13/viper) third party package.
- - Need implement  `func getViperConfig() viper_server.ViperConfig` function.

```go
package cache

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG Redis

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	PoolSize int    `mapstructure:"pool-size" json:"poolSize" yaml:"pool-size"`
}

// getViperConfig get initialize config
func getViperConfig() viper_server.ViperConfig {
	configName := "redis"
	db := fmt.Sprintf("%d", CONFIG.DB)
	poolSize := fmt.Sprintf("%d", CONFIG.PoolSize)
	return viper_server.ViperConfig{
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("deserialization data error: %v", err)
			}
			// config file change
			vi.SetConfigName(configName)
			vi.WatchConfig()
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("config file change:", e.Name)
				if err := vi.Unmarshal(&CONFIG); err != nil {
					fmt.Printf("deserialization data error: %v \n", err)
				}
			})
			return nil
		},
		// Note: When setting the default configuration value, there can be no other symbols such as spaces in front. It must be close to the left
		Default: []byte(`
db: ` + db + `
addr: "` + CONFIG.Addr + `"
password: "` + CONFIG.Password + `"
pool-size: ` + poolSize),
	}
}
```

- [zap_server]
- - Service logging.
- - Use [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap) third party package.
- - Through global variables `zap_server.ZAPLOG` record the log of the corresponding level.
```go
  zap_server.ZAPLOG.Info("Registration data table error", zap.Any("err", err))
  zap_server.ZAPLOG.Debug("Registration data table error", zap.Any("err", err))
  zap_server.ZAPLOG.Error("Registration data table error", zap.Any("err", err))
  ...
```

- [database]
- - database service [only support mysql now].
- - Use [gorm.io/gorm](https://github.com/go-gorm/gorm) third party package.
- - Through single instance `database.Instance()` operating data.
```go
  database.Instance().Model(&User{}).Where("name = ?","name").Find(&user)
  ...
```

- [casbin]
- - Access control management service.
- - Use [casbin](github.com/casbin/casbin/v2 ) third party package.
- - Through use `index.Use(casbin.Casbin())` middleware on route,implement interface authority authentication


- [cache]
- - Cache-driven service
- - Use [github.com/go-redis/redis](https://github.com/go-redis/redis) third party package.
- - Through single instance `cache.Instance()` operating data.

- [operation]
- - System operation log service.
- - Through use `index.Use(operation.OperationRecord())` middleware on route , realize the interface to automatically generate operation logs.

- [web]
- - web_iris Go-Iris web framework service.
- - Use [github.com/kataras/iris/v12](https://github.com/kataras/iris) third party package.
- - web framework service need implement `type WebFunc interface {}`  interface.
```go
// WebFunc web framework service interface
// - GetTestClient test client
// - GetTestLogin login for test
// - AddWebStatic add web static file
// - AddUploadStatic add upload file api
// - Run start program
type WebFunc interface {
	GetTestClient(t *testing.T) *httptest.Client
	GetTestLogin(t *testing.T, url string, res httptest.Responses, datas ...interface{}) *httptest.Client
	AddWebStatic(perfix string)
	AddUploadStatic()
	InitRouter() error
	Run()
}
```
#### Initialize database

##### Simple
- Use gorm's `AutoMigrate()` function to auto migrate database.
```go
package main

import (
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
  "github.com/snowlyg/iris-admin-rbac/iris/perm"
	"github.com/snowlyg/iris-admin-rbac/iris/role"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/operation"
)

func main() {
  	database.Instance().AutoMigrate(&perm.Permission{},&role.Role{},&user.User{},&operation.Oplog{})
}
```

##### Custom migrate tools
- Use `gormigrate` third party package. Tt's helpful for database migrate and program development.
- Detail is see  [iris-admin-cmd](https://github.com/snowlyg/iris-admin-example/blob/main/iris/cmd/main.go.

---

#### Getting started
- Get master package , Notice must use `master` version.
```sh
 go get github.com/snowlyg/iris-admin@master
```
- Add main.go file.
```go
package main

import (
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
)

func main() {
  wi := web_iris.Init()
	web.Start(wi)
}
```

#### Run project
- When you first run this cmd `go run main.go` , you can see some config files in  the `config` directory,
- and `rbac_model.conf` will be created in your project root directory.
```sh
go run main.go
```

#### Module
- You can use [iris-admin-rbac](https://github.com/snowlyg/iris-admin-rbac) package to add rbac function for your project quickly.
- Your can use AddModule() to add other modules .
```go
package main

import (
	rbac "github.com/snowlyg/iris-admin-rbac/iris"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
)

func main() {
	wi := web_iris.Init()
	rbacParty := web_iris.Party{
		Perfix:    "/api/v1",
		PartyFunc: rbac.Party(),
	}
	wi.AddModule(rbacParty)
	web.Start(web_iris.Init())
}
```

#### Default static file path
- A static file access path has been built in by default
- Static files will upload to `/static/upload` directory.
- You can set this config key `static-path` to change the default directory.
```yaml
system:
  addr: "127.0.0.1:8085"
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
- [iris](https://github.com/snowlyg/iris-admin-example/tree/main/iris)
- [gin](https://github.com/snowlyg/iris-admin-example/tree/main/gin)

#### RBAC
- [iris-admin-rbac](https://github.com/snowlyg/iris-admin-rbac)

#### Unit test and documentation
- Before start unit tests, you need to set two system environment variables `mysqlPwd` and `redisPwd`,that will be used when running the test instance。
- [helper/tests](https://github.com/snowlyg/helper/tree/main/tests) package the unit test used, it's  simple package base on [httpexpect/v2](https://github.com/gavv/httpexpect).
- [example for unit test](https://github.com/snowlyg/iris-admin-rbac/tree/main/iris/perm/tests)
- [example for unit test](https://github.com/snowlyg/iris-admin-rbac/tree/main/gin/authority/test)


#### Thanks

 - Thanks [JetBrains](https://www.jetbrains.com/?from=iris-admin)' supports .


*/
package doc
