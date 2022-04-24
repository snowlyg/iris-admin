<h1 align="center">IrisAdmin</h1>

<div align="center">
    <a href="https://codecov.io/gh/snowlyg/iris-admin"><img src="https://codecov.io/gh/snowlyg/iris-admin/branch/master/graph/badge.svg" alt="Code Coverage"></a>
    <a href="https://goreportcard.com/report/github.com/snowlyg/iris-admin"><img src="https://goreportcard.com/badge/github.com/snowlyg/iris-admin" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/snowlyg/iris-admin"><img src="https://godoc.org/github.com/snowlyg/iris-admin?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/snowlyg/iris-admin/blob/master/LICENSE"><img src="https://img.shields.io/github/license/snowlyg/iris-admin" alt="Licenses"></a>
</div>

简体中文 | [English](./README_EN.md)

#### 项目地址

[GITHUB](https://github.com/snowlyg/iris-admin) | [GITEE](https://gitee.com/snowlyg/iris-admin)

> 简单项目仅供学习，欢迎指点！

#### 相关文档

- [IRIS-ADMIN-DOC](https://doc.snowlyg.com)
- [IRIS V12 中文文档](https://github.com/snowlyg/iris/wiki)
- [godoc](https://pkg.go.dev/github.com/snowlyg/iris-admin?utm_source=godoc)

[![Gitter](https://badges.gitter.im/iris-go-tenancy/community.svg)](https://gitter.im/iris-go-tenancy/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) [![Join the chat at https://gitter.im/iris-go-tenancy/iris-admin](https://badges.gitter.im/iris-go-tenancy/iris-admin.svg)](https://gitter.im/iris-go-tenancy/iris-admin?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

#### iris 学习记录分享

- [Iris-go 项目登陆 API 构建细节实现过程](https://blog.snowlyg.com/iris-go-api-1/)

- [iris + casbin 从陌生到学会使用的过程](https://blog.snowlyg.com/iris-go-api-2/)

---

#### 简单使用

- 获取依赖包,注意必须带上 `master` 版本

```sh
 go get github.com/snowlyg/iris-admin@master
```

#### 项目介绍


##### 项目由多个插件构成,每个插件有不同的功能

- [viper_server]
  - 插件配置初始化,并生成本地配置文件
  - 使用 [github.com/spf13/viper](https://github.com/spf13/viper) 第三方包实现
  - 需要实现 `func getViperConfig() viper_server.ViperConfig` 方法

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

// getViperConfig 获取初始化配置
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
        return fmt.Errorf("反序列化错误: %v", err)
      }
      // 监控配置文件变化
      vi.SetConfigName(configName)
      vi.WatchConfig()
      vi.OnConfigChange(func(e fsnotify.Event) {
        fmt.Println("配置发生变化:", e.Name)
        if err := vi.Unmarshal(&CONFIG); err != nil {
          fmt.Printf("反序列化错误: %v \n", err)
        }
      })
      return nil
    },
    // 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
    Default: []byte(`
db: ` + db + `
addr: "` + CONFIG.Addr + `"
password: "` + CONFIG.Password + `"
pool-size: ` + poolSize),
  }
}
```

- [zap_server]
  - 插件日志记录
  - 使用 [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap) 第三方包实现
  - 通过全局变量 `zap_server.ZAPLOG` 记录对应级别的日志

```go
  zap_server.ZAPLOG.Info("注册数据表错误", zap.Any("err", err))
  zap_server.ZAPLOG.Debug("注册数据表错误", zap.Any("err", err))
  zap_server.ZAPLOG.Error("注册数据表错误", zap.Any("err", err))
  ...
```

- [database]
  - 数据插件 [目前仅支持 mysql]
  - 使用 [gorm.io/gorm](https://github.com/go-gorm/gorm) 第三方包实现
  - 通过单列 `database.Instance()` 操作数据

```go
  database.Instance().Model(&User{}).Where("name = ?","name").Find(&user)
  ...
```

- [casbin]
  - 权限控制管理插件
  - 使用 [casbin](github.com/casbin/casbin/v2 ) 第三方包实现
  - 并通过 `index.Use(casbin.Casbin())` 使用中间件,实现接口权限认证

- [cache]
  - 缓存驱动插件
  - 使用 [github.com/go-redis/redis](https://github.com/go-redis/redis) 第三方包实现
  - 通过单列 `cache.Instance()` 操作数据

- [operation]
  - 系统操作日志插件
  - 并通过 `index.Use(operation.OperationRecord())` 使用中间件,实现接口自动生成操作日志

- [cron_server]
  - 任务插件
  - 使用 [robfig/cron](https://github.com/robfig/cron) 第三方包实现
  - 通过单列 `cron_server.Instance()` 操作数据

```go
  cron_server.CronInstance().AddJob("@every 1m",YourJob)
  // 或者 
  cron_server.CronInstance().AddFunc("@every 1m",YourFunc)
  ...
```

- [web]
  - web_iris Go-Iris 框架插件
  - 使用 [github.com/kataras/iris/v12](https://github.com/kataras/iris) 第三方包实现
  - web 框架插件需要实现 `type WebFunc interface {}`  接口

```go
type WebBaseFunc interface {
  AddWebStatic(staticAbsPath, webPrefix string, paths ...string)
  AddUploadStatic(staticAbsPath, webPrefix string)
  InitRouter() error
  Run()
}

// WebFunc 框架插件接口
// - GetTestClient 测试客户端
// - GetTestLogin 测试登录
// - AddWebStatic 添加静态页面
// - AddUploadStatic 上传文件路径
// - Run 启动
type WebFunc interface {
  WebBaseFunc
}
```
  
---

#### 数据初始化

##### 简单初始化

- 使用原生方法 `AutoMigrate()` 自动迁移初始化数据表

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

##### 自定义迁移工具初始化

- 使用 `gormigrate` 第三方依赖包实现数据的迁移控制，方便后续的升级和开发
- 使用方法详情见 [iris-admin-cmd](https://github.com/snowlyg/iris-admin-example/blob/main/iris/cmd/main.go)

---

- 添加 main.go 文件

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

#### 启动项目

- 第一次启动项目后,配置文件会自动生成到 `config` 目录下.
- 同时会生成一个 `rbac_model.conf` 文件到项目根目录,该文件用于 casbin 权鉴的规则.

```sh
go run main.go
```

#### 添加模块

- 如果需要权鉴管理，可以使用 [iris-admin-rbac](https://github.com/snowlyg/iris-admin-rbac) 项目快速集成权鉴功能
- 可以使用 AddModule() 增加其他 admin模块

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

#### 设置静态文件路径

- 已经默认内置了一个静态文件访问路径
- 静态文件将会上传到 `/static/upload` 目录
- 可以修改配置项 `static-path` 修改默认目录

-

```yaml
system:
  addr: "127.0.0.1:8085"
  db-type: ""
  level: debug
  static-prefix: /upload
  time-format: "2006-01-02 15:04:05"
  web-prefix: /admin
  web-path: ./dist
```

#### 配合前端使用

- 编译前端页面默认 `dist` 目录
- 可以修改配置项 `web-path` 修改默认目录

-

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

- 前端页面参考/借用：【前端只简单实现预览效果】
- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin/tree/master/web)
- [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)

#### 简单用例

- [iris](https://github.com/snowlyg/iris-admin-example/tree/main/iris)
- [gin](https://github.com/snowlyg/iris-admin-example/tree/main/gin)

#### RBAC

- [iris-admin-rbac](https://github.com/snowlyg/iris-admin-rbac)

#### 接口单元测试和接口文档

- 测试前,需要设置 `mysqlPwd`, `mysqlAddr` 系统环境变量，运行测试实例的时候将会使用到它们。
- 测试使用依赖库 [helper/tests](https://github.com/snowlyg/helper/tree/main/tests) 是基于 [httpexpect/v2](https://github.com/gavv/httpexpect) 的简单封装
- [接口单元测试例子](https://github.com/snowlyg/iris-admin-rbac/tree/main/iris/perm/tests)
- [接口单元测试例子](https://github.com/snowlyg/iris-admin-rbac/tree/main/gin/authority/test)

接口单元测试需要新建 `main_test.go` 文件,该文件定义了单元测试的一些通用基础步骤：
***建议采用docker部署mysql,否则测试失败会有大量测试数据库遗留***
1.测试数据库的数据库的创建和摧毁（每个单元测试都会新建不同的数据库，以隔离数据对单元测试结果的影响）
2.数据表的新建和表数据的填充
3. `PartyFunc` , `SeedFunc` 方法需要根据对应的测试模块自定义
内容如下所示:

***main_test.go***

```go
package test

import (
  "os"
  "testing"

  "github.com/snowlyg/httptest"
  rbac "github.com/snowlyg/iris-admin-rbac/gin"
  "github.com/snowlyg/iris-admin/server/web/common"
  "github.com/snowlyg/iris-admin/server/web/web_gin"
)

var TestServer *web_gin.WebServer
var TestClient *httptest.Client

func TestMain(m *testing.M) {

  var uuid string
  uuid, TestServer = common.BeforeTestMainGin(rbac.PartyFunc, rbac.SeedFunc)
  code := m.Run()
  common.AfterTestMain(uuid, true)

  os.Exit(code)
}

```

***index_test.go***

```go
package test

import (
  "fmt"
  "net/http"
  "path/filepath"
  "testing"

  "github.com/snowlyg/helper/str"
  "github.com/snowlyg/httptest"
  rbac "github.com/snowlyg/iris-admin-rbac/gin"
  "github.com/snowlyg/iris-admin/g"
  "github.com/snowlyg/iris-admin/server/web"
  "github.com/snowlyg/iris-admin/server/web/web_gin/response"
)

var (
  url = "/api/v1/admin"
)

func TestList(t *testing.T) {
  TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
  TestClient.Login(rbac.LoginUrl, nil)
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
          {Key: "id", Value: 1, Type: "ge"},
          {Key: "nickName", Value: "超级管理员"},
          {Key: "username", Value: "admin"},
          {Key: "headerImg", Value: "http://xxxx/head.png"},
          {Key: "status", Value: g.StatusTrue},
          {Key: "isShow", Value: g.StatusFalse},
          {Key: "phone", Value: "13800138000"},
          {Key: "email", Value: "admin@admin.com"},
          {Key: "authorities", Value: []string{"超级管理员"}},
          {Key: "updatedAt", Value: "", Type: "notempty"},
          {Key: "createdAt", Value: "", Type: "notempty"},
        },
      }},
      {Key: "total", Value: 0, Type: "ge"},
    }},
  }
  TestClient.GET(fmt.Sprintf("%s/getAll", url), pageKeys, httptest.RequestParams)
}

func TestCreate(t *testing.T) {
  TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
  TestClient.Login(rbac.LoginUrl, nil)
  if TestClient == nil {
    return
  }

  data := map[string]interface{}{
    "nickName":     "测试名称",
    "username":     "create_test_username",
    "authorityIds": []uint{web.AdminAuthorityId},
    "email":        "get@admin.com",
    "phone":        "13800138001",
    "password":     "123456",
  }
  id := Create(TestClient, data)
  if id == 0 {
    t.Fatalf("测试添加用户失败 id=%d", id)
  }
  defer Delete(TestClient, id)
}

func TestUpdate(t *testing.T) {

  TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
  TestClient.Login(rbac.LoginUrl, nil)
  if TestClient == nil {
    return
  }
  data := map[string]interface{}{
    "nickName":     "测试名称",
    "username":     "create_test_username_for_update",
    "authorityIds": []uint{web.AdminAuthorityId},
    "email":        "get@admin.com",
    "phone":        "13800138001",
    "password":     "123456",
  }
  id := Create(TestClient, data)
  if id == 0 {
    t.Fatalf("测试添加用户失败 id=%d", id)
  }
  defer Delete(TestClient, id)

  update := map[string]interface{}{
    "nickName": "测试名称",
    "email":    "get@admin.com",
    "phone":    "13800138003",
    "password": "123456",
  }

  pageKeys := httptest.Responses{
    {Key: "status", Value: http.StatusOK},
    {Key: "message", Value: response.ResponseOkMessage},
  }
  TestClient.PUT(fmt.Sprintf("%s/updateAdmin/%d", url, id), pageKeys, update)
}

func TestGetById(t *testing.T) {
  TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
  TestClient.Login(rbac.LoginUrl, nil)
  if TestClient == nil {
    return
  }
  data := map[string]interface{}{
    "nickName":     "测试名称",
    "username":     "create_test_username_for_get",
    "email":        "get@admin.com",
    "phone":        "13800138001",
    "authorityIds": []uint{web.AdminAuthorityId},
    "password":     "123456",
  }
  id := Create(TestClient, data)
  if id == 0 {
    t.Fatalf("测试添加用户失败 id=%d", id)
  }
  defer Delete(TestClient, id)
  pageKeys := httptest.Responses{
    {Key: "status", Value: http.StatusOK},
    {Key: "message", Value: response.ResponseOkMessage},
    {Key: "data", Value: httptest.Responses{
      {Key: "id", Value: 1, Type: "ge"},
      {Key: "nickName", Value: data["nickName"].(string)},
      {Key: "username", Value: data["username"].(string)},
      {Key: "status", Value: g.StatusTrue},
      {Key: "email", Value: data["email"].(string)},
      {Key: "phone", Value: data["phone"].(string)},
      {Key: "isShow", Value: g.StatusTrue},
      {Key: "headerImg", Value: "http://xxxx/head.png"},
      {Key: "updatedAt", Value: "", Type: "notempty"},
      {Key: "createdAt", Value: "", Type: "notempty"},
      {Key: "createdAt", Value: "", Type: "notempty"},
      {Key: "authorities", Value: []string{"超级管理员"}},
    },
    },
  }
  TestClient.GET(fmt.Sprintf("%s/getAdmin/%d", url, id), pageKeys)
}

func TestChangeAvatar(t *testing.T) {
  TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
  TestClient.Login(rbac.LoginUrl, nil)
  if TestClient == nil {
    return
  }
  data := map[string]interface{}{
    "headerImg": "/avatar.png",
  }
  pageKeys := httptest.Responses{
    {Key: "status", Value: http.StatusOK},
    {Key: "message", Value: response.ResponseOkMessage},
  }
  TestClient.POST(fmt.Sprintf("%s/changeAvatar", url), pageKeys, data)

  profile := httptest.Responses{
    {Key: "status", Value: http.StatusOK},
    {Key: "message", Value: response.ResponseOkMessage},
    {Key: "data", Value: httptest.Responses{
      {Key: "id", Value: 1, Type: "ge"},
      {Key: "nickName", Value: "超级管理员"},
      {Key: "username", Value: "admin"},
      {Key: "headerImg", Value: filepath.ToSlash(web.ToStaticUrl("/avatar.png"))},
      {Key: "status", Value: g.StatusTrue},
      {Key: "isShow", Value: g.StatusFalse},
      {Key: "phone", Value: "13800138000"},
      {Key: "email", Value: "admin@admin.com"},
      {Key: "authorities", Value: []string{"超级管理员"}},
      {Key: "updatedAt", Value: "", Type: "notempty"},
      {Key: "createdAt", Value: "", Type: "notempty"},
    },
    },
  }
  TestClient.GET(fmt.Sprintf("%s/profile", url), profile)
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
  return TestClient.POST(fmt.Sprintf("%s/createAdmin", url), pageKeys, data).GetId()
}

func Delete(TestClient *httptest.Client, id uint) {
  pageKeys := httptest.Responses{
    {Key: "status", Value: http.StatusOK},
    {Key: "message", Value: response.ResponseOkMessage},
  }
  TestClient.DELETE(fmt.Sprintf("%s/deleteAdmin/%d", url, id), pageKeys)
}

```

#### 感谢

[JetBrains](https://www.jetbrains.com/?from=iris-admin) 对本项目的支持。

#### 打赏
> 您的打赏将用于支付网站运行，会在项目介绍中特别鸣谢您
- [为爱发电](https://afdian.net/@snowlyg/plan)
- [donating](https://paypal.me/snowlyg?country.x=C2&locale.x=zh_XC)
