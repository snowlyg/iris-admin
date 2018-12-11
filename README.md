# iris golang for api

#### 项目介绍
- 采用 iris 框架目后台 api （来自公司内部项目重构）
- 采用了 gorm 数据库模块 和 jwt 的单点登陆认证方式
- 测试默认使用了 sqlite3 数据库
- 修改了项目文件结构，重新采用了 models , controllers 的结构（上次用的那种方式实在是看不下去了）。

---

#### 项目初始化

>拉取项目 
```
git clone https://github.com/569616226/IrisApiProject.gits
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


国内的 cdn


https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```

>访问文档，从浏览器直接打开 apidoc/index.html 文件

---


##### 参考资料
- [zuoyanart/pizzaCmsApi](https://github.com/zuoyanart/pizzaCmsApi) 
- [beego应用做纯API后端如何使用jwt实现无状态权限验证](https://www.cnblogs.com/lrj567/p/6209872.html)
- [Go实现jwt](https://blog.csdn.net/zxy_666/article/details/80021331)
- [Go实战--golang中使用JWT(JSON Web Token)](https://blog.csdn.net/wangshubo1989/article/details/74529333)