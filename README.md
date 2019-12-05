# IrisApi
###### Iris + Vue + mysql + redis + jwt

#### 项目介绍
- 采用 iris 框架目后台api [IrisApiProject](https://github.com/snowlyg/IrisApiProject.gits)
- 采用了 gorm 数据库模块 和 jwt 的单点登陆认证方式
- 测试默认使用了 sqlite3 数据库
---

```
项目更新到了 iris v12,对应的也要 iris 升级 
 go get -u github.com/kataras/iris

如果要用旧版本的 iris ,需要克隆本项目 1.0 版本
```

#### 项目目录结构
- apidoc 接口文档目录
- caches redis缓存目录
- config 项目配置文件目录
- controllers 控制器文件目录
- database 数据库文件目录
- middleware 中间件文件目录
- models 模型文件目录
- routes 路由文件
- resources 前端文件
- tmp 测试数据库 sqlite3 文件目录
- tools 其他公用方法目录
---

#### api项目初始化

>拉取项目

```
git clone https://github.com/snowlyg/IrisApiProject.git
```
//github 太慢可以用 gitee
```
git clone https://gitee.com/dtouyu/IrisApiProject.git
```

>加载依赖管理包

使用 [gopm](https://gopm.io/) 管理包
```

go get -v -u github.com/gpmgo/gopm

# 查看当前工程依赖
gopm list
# 显示依赖详细信息
gopm list -v
# 列出文件依赖
gopm list -t [file]
# 拉取依赖到缓存目录
gopm get -r xxx
# 仅下载当前指定的包
gopm get -d xxx
# 拉取依赖到$GOPATH
gopm get -g xxx
# 检查更新所有包
gopm get -u xxx
# 拉取到当前所在目录
gopm get -l xxx
# 运行当前目录程序
gopm run
# 生成当前工程的 gopmfile 文件用于包管理
gopm gen -v
# 根据当前项目 gopmfile 链接依赖并执行 go install
gopm install -v
# 更新当前依赖
gopm update -v
# 清理临时文件
gopm clean
# 编译到当前目录
gopm bin

```

>项目配置文件 /config/config.toml

```
cp config.toml.example config.toml
```

>运行项目 

[gowatch](https://gitee.com/silenceper/gowatch)
```
go get github.com/silenceper/gowatch

gowatch //安装 gowatch 后才可以使用这个命令，不然只能使用

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
>打开 /resource/apiDoc/index.html 修改里面的

```
https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js

国内的 cdn


https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```

>访问文档，从浏览器直接打开 http://localhost/apiDoc

---

#### 登录项目
输入地址 http://localhost:80

//在 conig/config.toml 内配置 `LoginUserName` 和 `LoginPwd`

项目管理员账号 ： username
项目管理员密码 ： password

##### 问题总结

[问题总结](ERRORS.MD)


##### 参考资料
- [zuoyanart/pizzaCmsApi](https://github.com/zuoyanart/pizzaCmsApi) 
- [beego应用做纯API后端如何使用jwt实现无状态权限验证](https://www.cnblogs.com/lrj567/p/6209872.html)
- [Go实现jwt](https://blog.csdn.net/zxy_666/article/details/80021331)
- [Go实战--golang中使用JWT(JSON Web Token)](https://blog.csdn.net/wangshubo1989/article/details/74529333)
