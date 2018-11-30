# iris golang for api

#### 项目介绍
- 采用 iris 框架目后台 api （来自公司内部项目重构）
- 采用了 gorm 数据库模块 和 jwt 的单点登陆认证方式
- 测试默认使用了 sqlite3 数据库
- 最开始想用 mvc 将 controller router model 分离开，还是受到 laravel 影响比较大。
  不过在测试的时候不好切换数据库，使用了配置文件去标记运行环境。还是每次都要手动切换，十分麻烦。
  如果不小心忘记改了，还会清空正常数据库的数据。
  最后把 router model controller 都放在了 main 包下，用 c_ m_ l_ 的前缀方式去区分。 

---

#### 项目初始化

>拉取项目 
```
git clone https://git.dev.tencent.com/Dreamfish/IrisYouQiKangApi.git
```

>加载依赖管理包
```
 本来是用 godep 管理的，使用后发现还是是有问题。暂时不使用依赖管理包，依赖要自行下载。
```

>项目配置文件 /config/config.toml

```

cp config.toml.example config.toml
```

>运行项目 

```
gowatch //安装了 gowatch 热加载工具可以直接运行

go run main.go // go 命令
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