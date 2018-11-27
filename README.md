# iris golang for api

#### 项目介绍
- 采用 iris 框架目后台 api （来自公司内部项目重构）
- 采用了 gorm 数据库模块 和 jwt 的单点登陆认证方式
- 自己封装下一下 http 测试功能 （考虑重构成插件）
- 测试默认使用了 sqlite 数据库
- 使用 godep 依赖管理包（对比了其他包，还是选择用这个）

---

#### 项目初始化

>拉取项目 
```
git clone https://git.dev.tencent.com/Dreamfish/IrisYouQiKangApi.git
```

>加载 godep 依赖管理包（已经安装可以跳过）
```
go get -u -v github.com/tools/godep
```


>加载 go 依赖库 
```
godep restore

```


>项目配置文件 config.toml

```
[hotreload]
  suffixes = [".go"]
  ignore = []

[app]

  #app 环境 test local
  env = "testing"
  #app 名称
  name = "IrisYouQiKang"
  #app url
  url  = "http://localhost"
  #app 文档地址
  doc = "./apidoc"
  addr = ":80"

  [app.logger]
    level = "INFO"
    name = "application"

#数据库驱动
[database]
    dirver = "mysql"

[mysql]
  connect = "root:you_password@/you_database_name?charset=utf8&parseTime=True&loc=Local"

[mongodb]
  connect = "mongodb://root:123456@127.0.0.1:27017/admin"

[sqlite]
  connect = "/tmp/gorm.db"

#reids
[redis]
  Addr = "127.0.0.1:6379"
  Password = ""
  DB = 0

[neo4j]
    connect = "http://10.10.43.111:7474/db/data"

[test]
    #测试登陆用户名
    LoginUserName = "you_test_user_name"
    #测试用户名
    LoginName = "you_test_name"
    #测试用户密码
    LoginPwd = "you_test_user_password"
    #测试数据库驱动
    DataBaseDriver = "sqlite3"
    #测试数据库
    DataBaseConnect = "/tmp/gorm.db"

```

---
##### 单元测试 
>单元测试我做了简单的封装，可以使用下面的方法写 get 请求的测试
>
>将测试文件放在项目目录下，执行  `godep save` 命令。会出现没有使用的依赖库，却报错依赖库不存在的问题。
```
godep: Package (github.com/sergi/go-diff/diffmatchpatch) not found
```
>
>将单元测试的代码 移动到 `tests` 文件夹下， 可以解决执行  `godep save` 命令出现没有使用的依赖库没有加载的报错。


```
//设置测试数据表
	SetTestTableName("users")

	//创建系统管理员，测试 users 表需要手动创建。
	//其他模型测试不需要手动创建
	aul := CreaterSystemAdmin()
	users := []*models.AdminUserTranform{aul}

	//发起 http 请求
	//Url        string      //测试路由
	//Object     interface{} //发送的json 对象
	//StatusCode int         //返回的 http 状态码
	//Status     bool        //返回的状态
	//Msg        string      //返回提示信息
	//Data       interface{} //返回数据
	bc := BaseCase{"/v1/admin/users", nil, iris.StatusOK, true, "操作成功", users}
	bc.get(t)
```

> post 请求

```
//输入错误的登陆密码
func TestUserLoginWithErrorPwd(t *testing.T) {
	oj := map[string]interface{}{
		"username": system.Config.Get("test.LoginUserName").(string),
		"password": "admin",
	}
	bc := BaseCase{"/v1/admin/login", oj, iris.StatusOK, false, "用户名或密码错误", nil}
	bc.post(t)
}
```



---

##### api 文档使用
自动生成文档 (访问过接口就会自动成功)
因为原生的 jquery.min.js 里面的 cdn 是使用国外的，访问很慢。
有条件的可以开个 vpn ,如果没有可以根据下面的方法修改一下，访问就很快了
>打开 apidoc/index.html 修改里面的

```
https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js


为国内的 cdn


https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```

>访问文档，从浏览器直接打开 apidoc/index.html 文件


---


##### 参考资料
- [zuoyanart/pizzaCmsApi](https://github.com/zuoyanart/pizzaCmsApi) 
- [beego应用做纯API后端如何使用jwt实现无状态权限验证](https://www.cnblogs.com/lrj567/p/6209872.html)
- [Go实现jwt](https://blog.csdn.net/zxy_666/article/details/80021331)
- [Go实战--golang中使用JWT(JSON Web Token)](https://blog.csdn.net/wangshubo1989/article/details/74529333)