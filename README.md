# iris golang for youqikangapi

采用 iris 框架重构优企康项目后台 api
采用了 xorm 和 jwt 的登陆认证方式

参考了 [zuoyanart/pizzaCmsApi](https://github.com/zuoyanart/pizzaCmsApi) 的项目

项目配置文件 config.toml

```
#
# settings related to recompiling and reruning app when source changes
#
[hotreload]
  # file suffixes to watch for changes
  suffixes = [".go"]

  # files/directories to ignore
  ignore = []

#
# general application settings
#
[app]
  # additional application arguments
  args = []
  addr = ":3000"

  [app.logger]
    level = "INFO"
    name = "application"

# 运行参数
[mysql]
  connect = "root:UHC0JC5s6DEg9BRXYuDJnqbdl1ecL4gV@/goyouqikang?charset=utf8&parseTime=True&loc=Local"
  MaxIdle = 10
  MaxOpen = 100
[mongodb]
  connect = "mongodb://root:123456@127.0.0.1:27017/admin"
[redis]
  connect = "127.0.0.1:6379"
  db = 0
  maxidle = 100
  maxactive = 1000
[neo4j]
    connect = "http://10.10.43.111:7474/db/data"
```



自动成功 api 文档 (访问过接口就会自动成功)
1.打开 apidoc.html 修改里面的

```
https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js
```

为国内的 cdn

```
https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```
