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
- [IRIS-ADMIN-DOC](https://doc.snowlyg.com)
- [IRIS V12 document for chinese](https://github.com/snowlyg/iris/wiki)
- [godoc](https://pkg.go.dev/github.com/snowlyg/iris-admin?utm_source=godoc)


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

- [cron_server]
- - Job server
- - Use [robfig/cron](https://github.com/robfig/cron) third party package.
- - Through single instance `cron_server.Instance()` to add job or func.
```go
  cron_server.CronInstance().AddJob("@every 1m",YourJob)
  // 或者 
  cron_server.CronInstance().AddFunc("@every 1m",YourFunc)
  ...
```

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
  db-type: ""
  level: debug
  static-prefix: /upload
  time-format: "2006-01-02 15:04:05"
  web-path: ./dist
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
  wi.AddUploadStatic("/upload", "/var/static")
  wi.AddWebStatic("/", "/var/static")
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

接口单元测试需要新建 `main_test.go` 文件,该文件定义了单元测试的一些通用基础步骤：
***建议采用docker部署mysql,否则测试失败会有大量测试数据库遗留***
1.测试数据库的数据库的创建和摧毁（每个单元测试都会新建不同的数据库，以隔离数据对单元测试结果的影响）
2.数据表的新建和表数据的填充
3. `PartyFunc` , `SeedFunc` 方法需要根据对应的测试模块自定义
内容如下所示:
```go
package test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/snowlyg/httptest"
	rbac "github.com/snowlyg/iris-admin-rbac/gin"
	web_tests "github.com/snowlyg/iris-admin/server/web/tests"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
)

var TestServer *web_gin.WebServer
var TestClient *httptest.Client

func TestMain(m *testing.M) {
	var uuid string
	uuid, TestServer = web_tests.BeforeTestMainGin(rbac.PartyFunc, rbac.SeedFunc)
	code := m.Run()
	web_tests.AfterTestMain(uuid, true)
	os.Exit(code)
}
```

4.然后添加对应的单元测试
```go
package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/snowlyg/httptest"
	rbac "github.com/snowlyg/iris-admin-rbac/gin"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
	"github.com/snowlyg/multi"
)

var (
	url = "/api/v1/authority" // url
)

func TestList(t *testing.T) {
	if TestServer == nil {
		t.Error("测试服务初始化失败")
		return
	}

	TestClient = TestServer.GetTestLogin(t, rbac.LoginUrl, rbac.LoginResponse)
	if TestClient == nil {
		return
	}
	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: httptest.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "list", Value: []httptest.Responses{
				{
					{Key: "id", Value: web.DeviceAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "设备用户"},
					{Key: "authorityType", Value: multi.GeneralAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
				{
					{Key: "id", Value: web.LiteAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "小程序用户"},
					{Key: "authorityType", Value: multi.GeneralAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
				{
					{Key: "id", Value: web.TenancyAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "商户管理员"},
					{Key: "authorityType", Value: multi.TenancyAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
				{
					{Key: "id", Value: web.AdminAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "超级管理员"},
					{Key: "authorityType", Value: multi.AdminAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
			}},
			{Key: "total", Value: 0, Type: "ge"},
		}},
	}
	requestParams := map[string]interface{}{"page": 1, "pageSize": 10, "orderBy": "id"}
	TestClient.GET(fmt.Sprintf("%s/getAuthorityList", url), pageKeys, requestParams)
}

func TestGetAdminAuthorityList(t *testing.T) {
	if TestServer == nil {
		t.Error("测试服务初始化失败")
		return
	}

	TestClient = TestServer.GetTestLogin(t, rbac.LoginUrl, rbac.LoginResponse)
	if TestClient == nil {
		return
	}
	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: httptest.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "list", Value: []httptest.Responses{
				{
					{Key: "id", Value: web.AdminAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "超级管理员"},
					{Key: "authorityType", Value: multi.AdminAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
			}},
			{Key: "total", Value: 0, Type: "ge"},
		}},
	}
	requestParams := map[string]interface{}{"page": 1, "pageSize": 10, "orderBy": "id"}
	TestClient.GET(fmt.Sprintf("%s/getAdminAuthorityList", url), pageKeys, requestParams)
}

func TestGetTenancyAuthorityList(t *testing.T) {
	if TestServer == nil {
		t.Error("测试服务初始化失败")
		return
	}

	TestClient = TestServer.GetTestLogin(t, rbac.LoginUrl, rbac.LoginResponse)
	if TestClient == nil {
		return
	}
	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: httptest.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "list", Value: []httptest.Responses{
				{
					{Key: "id", Value: web.TenancyAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "商户管理员"},
					{Key: "authorityType", Value: multi.TenancyAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
			}},
			{Key: "total", Value: 0, Type: "ge"},
		}},
	}
	requestParams := map[string]interface{}{"page": 1, "pageSize": 10, "orderBy": "id"}
	TestClient.GET(fmt.Sprintf("%s/getTenancyAuthorityList", url), pageKeys, requestParams)
}

func TestGetGeneralAuthorityList(t *testing.T) {
	if TestServer == nil {
		t.Error("测试服务初始化失败")
		return
	}

	TestClient = TestServer.GetTestLogin(t, rbac.LoginUrl, rbac.LoginResponse)
	if TestClient == nil {
		return
	}
	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: httptest.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "list", Value: []httptest.Responses{
				{
					{Key: "id", Value: web.DeviceAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "设备用户"},
					{Key: "authorityType", Value: multi.GeneralAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
				{
					{Key: "id", Value: web.LiteAuthorityId, Type: "ge"},
					{Key: "authorityName", Value: "小程序用户"},
					{Key: "authorityType", Value: multi.GeneralAuthority},
					{Key: "parentId", Value: 0},
					{Key: "defaultRouter", Value: "dashboard"},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
			}},
			{Key: "total", Value: 0, Type: "ge"},
		}},
	}
	requestParams := map[string]interface{}{"page": 1, "pageSize": 10, "orderBy": "id"}
	TestClient.GET(fmt.Sprintf("%s/getGeneralAuthorityList", url), pageKeys, requestParams)
}

func TestCreate(t *testing.T) {
	if TestServer == nil {
		t.Error("测试服务初始化失败")
		return
	}

	TestClient = TestServer.GetTestLogin(t, rbac.LoginUrl, rbac.LoginResponse)
	if TestClient == nil {
		return
	}
	data := map[string]interface{}{
		"authorityName": "test_authorityName_for_create",
		"parentId":      0,
		"authorityType": multi.AdminAuthority,
	}
	id := Create(TestClient, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(TestClient, id)
}

func TestUpdate(t *testing.T) {
	if TestServer == nil {
		t.Error("测试服务初始化失败")
		return
	}

	TestClient = TestServer.GetTestLogin(t, rbac.LoginUrl, rbac.LoginResponse)
	if TestClient == nil {
		return
	}
	data := map[string]interface{}{
		"authorityName": "test_authorityName_for_update",
		"parentId":      0,
		"authorityType": multi.AdminAuthority,
	}
	id := Create(TestClient, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(TestClient, id)

	update := map[string]interface{}{
		"authorityName": "test_authorityName_for_update1",
		"parentId":      0,
		"authorityType": multi.AdminAuthority,
	}

	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	TestClient.PUT(fmt.Sprintf("%s/updateAuthority/%d", url, id), pageKeys, update)
}

func TestCopyAuthority(t *testing.T) {
	if TestServer == nil {
		t.Error("测试服务初始化失败")
		return
	}

	TestClient = TestServer.GetTestLogin(t, rbac.LoginUrl, rbac.LoginResponse)
	if TestClient == nil {
		return
	}
	data := map[string]interface{}{
		"authorityName": "test_authorityName_for_copy",
		"parentId":      0,
		"authorityType": multi.AdminAuthority,
	}
	id := Create(TestClient, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(TestClient, id)

	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	copy := map[string]interface{}{
		"authorityName": "test_authorityName_after_copy",
	}

	TestClient.POST(fmt.Sprintf("%s/copyAuthority/%d", url, id), pageKeys, copy)
}

func Create(TestClient *httptest.Client, data map[string]interface{}) uint {
	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: httptest.Responses{
			{Key: "id", Value: 1, Type: "ge"},
		},
		},
	}
	return TestClient.POST(fmt.Sprintf("%s/createAuthority", url), pageKeys, data).GetId()
}

func Delete(TestClient *httptest.Client, id uint) {
	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	TestClient.DELETE(fmt.Sprintf("%s/deleteAuthority/%d", url, id), pageKeys)
}
```

#### Thanks 

 - Thanks [JetBrains](https://www.jetbrains.com/?from=iris-admin)' supports .


