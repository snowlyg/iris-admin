# 问题总结

> 运行测试 go test ,报错 exec: "gcc": executable file not found in %PATH%
- 系统没有安装 gcc ,或者 gcc 没有配置到系统环境变量
- 安装 gcc 并设置相应系统变量

> 安装依赖包，出现无法安装或者安装太慢（国内访问都会有限制）
 ```
   有条件的可以使用 vpn , 或者到 github 去下载对应的依赖包文件，放到对应的加载路径。
   github 如果下载或者 clone 还是很慢，可以用 [gopm.io](https://gopm.io/download) 加速下载
```


> 启动报错如下内容：
 ```
   GetPermissionByNameError:record not found创建账号权限
   CreatePermissionError:Error 1366: Incorrect string value: '\xE5\x88\x9B\xE5\xBB\xBA...' for column 'display_name' at row 1GetRoleByNameErr:record not found创建角色
   
   (/home/coding/go/src/IrisApiProject/models/roles.go:108) 
   [2019-09-27 12:02:22]  Error 1366: Incorrect string value: '\xE8\xB6\x85\xE7\xBA\xA7...' for column 'display_name' at row 1 
   CreateRoleErr:Error 1366: Incorrect string value: '\xE8\xB6\x85\xE7\xBA\xA7...' for column 'display_name' at row 1
   (/home/coding/go/src/IrisApiProject/models/permission.go:99) 
   [2019-09-27 12:02:22]  Error 1366: Incorrect string value: '\xE5\x88\x9B\xE5\xBB\xBA...' for column 'display_name' at row 1 
```
- 数据库字符格式不对，无法识别中文
- 将数据库字符设置 utf8 字符集

> 登录显示前端报错 405 Not Allowed：
 ```
默认端口为 8081，如果修改了端口。需要下载 [IrisApiVueAdmin](https://github.com/snowlyg/IrisApiVueAdmin) 项目修改端口后覆盖前端代码
