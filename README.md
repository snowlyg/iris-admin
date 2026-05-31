<h1 align="center">IrisAdmin</h1>

[![Build Status](https://app.travis-ci.com/snowlyg/iris-admin.svg?branch=master)](https://app.travis-ci.com/snowlyg/iris-admin)
[![LICENSE](https://img.shields.io/github/license/snowlyg/iris-admin)](https://github.com/snowlyg/iris-admin/blob/master/LICENSE)
[![go doc](https://godoc.org/github.com/snowlyg/iris-admin?status.svg)](https://godoc.org/github.com/snowlyg/iris-admin)
[![go report](https://goreportcard.com/badge/github.com/snowlyg/iris-admin)](https://goreportcard.com/badge/github.com/snowlyg/iris-admin)
[![Build Status](https://codecov.io/gh/snowlyg/iris-admin/branch/master/graph/badge.svg)](https://codecov.io/gh/snowlyg/iris-admin)

简体中文 | [English](./README_EN.md)

#### 项目地址

[GITHUB](https://github.com/snowlyg/iris-admin)

IrisAdmin 是一个 Apache-2.0 开源的 Go Web Admin / RBAC 脚手架项目，围绕 Iris/Gin、GORM、Casbin、Redis、Docker、JWT 鉴权和 REST API 工作流沉淀通用后台服务实践。

项目由 [snowlyg](https://github.com/snowlyg) 维护，欢迎通过 issue 和 pull request 参与文档、测试、依赖升级和安全改进。

#### 相关文档

- [IRIS-ADMIN-DOC](https://doc.snowlyg.com)
- [IRIS V12 中文文档](https://github.com/snowlyg/iris/wiki)
- [godoc](https://pkg.go.dev/github.com/snowlyg/iris-admin?utm_source=godoc)
- [2026 Roadmap](./ROADMAP.md)
- [Contributing](./CONTRIBUTING.md)
- [Security Policy](./SECURITY.md)

<a href="https://gitter.im/iris-go-tenancy/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge"><img src="https://badges.gitter.im/iris-go-tenancy/community.svg" alt="e9939a7e92f32337871feb22e06bd05a.jpeg" border="0" width=120 /></a>
<a href="https://discord.gg/pytCGMSBgA"> <img src="https://www.svgrepo.com/show/353655/discord-icon.svg" alt="e9939a7e92f32337871feb22e06bd05a.jpeg" border="0" width=30 /></a>

#### iris 学习记录分享

- [Iris-go 项目登陆 API 构建细节实现过程](https://snowlyg.github.io/posts/iris-go-api-1/)

- [iris + casbin 从陌生到学会使用的过程](https://snowlyg.github.io/posts/iris-go-api-2/)

---

#### 简单使用

- 获取依赖包,注意必须带上 `master` 版本

```sh
 go get github.com/snowlyg/iris-admin@master
```

#### Maintainer status

IrisAdmin 作为长期维护的开源项目，后续维护重点包括：

- 升级 Go 依赖和 CI 工作流。
- 补强 RBAC / Casbin / JWT / middleware 等关键路径测试。
- 改进 Docker 和本地 quickstart 文档。
- 完善英文文档，方便更多 Go 开发者阅读和参与。
- 审查鉴权、权限校验和 API 中间件等安全敏感代码路径。

#### 打赏

> 您的打赏将用于支付网站运行，会在项目介绍中特别鸣谢您
- [爱发电](https://afdian.net/@snowlyg/plan)
- [donating](https://paypal.me/snowlyg?country.x=C2&locale.x=zh_XC)
