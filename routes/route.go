package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
	"github.com/snowlyg/blog/controllers"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/middleware"
	"github.com/snowlyg/easygorm"
	"path/filepath"
	"time"
)

const maxSize = 5 << 20 // 5MB

func App(api *iris.Application) {
	api.UseRouter(middleware.CrsAuth())
	app := api.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.HandleDir("/uploads", iris.Dir(filepath.Join(libs.CWD(), "uploads")))

		v1 := app.Party("/v1")
		{
			// 是否开启接口请求频率限制
			if !libs.Config.Limit.Disable {
				limitV1 := rate.Limit(libs.Config.Limit.Limit, libs.Config.Limit.Burst, rate.PurgeEvery(time.Minute, 5*time.Minute))
				v1.Use(limitV1)
			}
			v1.Get("/profile", controllers.GetAdminInfo)
			v1.PartyFunc("/article", func(article iris.Party) {
				article.Get("/", controllers.GetAllPublishedArticles)
				article.Get("/{id:uint}", controllers.GetPublishedArticle)
				article.Get("/like/{id:uint}", controllers.GetPublishedArticleLike)
			})
			v1.PartyFunc("/config/{key:string}", func(configs iris.Party) {
				configs.Get("/", controllers.GetConfig)
			})
			v1.PartyFunc("/types", func(articleType iris.Party) {
				articleType.Get("/", controllers.GetAllTypes)
			})
			v1.PartyFunc("/tags", func(articleTag iris.Party) {
				articleTag.Get("/", controllers.GetAllTags)
			})
			v1.PartyFunc("/docs", func(docs iris.Party) {
				docs.Get("/", controllers.GetAllDocs)
				docs.Get("/{id:uint}", controllers.GetDoc)
			})
			v1.PartyFunc("/chapter", func(chapter iris.Party) {
				chapter.Get("/", controllers.GetAllPublishedChapters)
				chapter.Get("/{id:uint}", controllers.GetPublishedChapter)
				chapter.Get("/like/{id:uint}", controllers.GetPublishedChapterLike)
			})
			v1.Post("/admin/login", controllers.UserLogin)
			v1.PartyFunc("/admin", func(admin iris.Party) {
				casbinMiddleware := middleware.New(easygorm.Egm.Enforcer)            //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
				admin.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
				admin.Post("/logout", controllers.UserLogout).Name = "退出"
				admin.Get("/expire", controllers.UserExpire).Name = "刷新 token"
				admin.Get("/profile", controllers.GetProfile).Name = "个人信息"
				admin.Put("/change_avatar", controllers.ChangeAvatar).Name = "修改头像"
				admin.Get("/dashboard", controllers.Dashboard).Name = "数据统计"
				admin.Post("/upload_file", iris.LimitRequestBodySize(maxSize+1<<20), controllers.UploadFile).Name = "上传文件"
				admin.PartyFunc("/article", func(aritcle iris.Party) {
					aritcle.Get("/", controllers.GetAllArticles).Name = "文章列表"
					aritcle.Get("/{id:uint}", controllers.GetArticle).Name = "文章详情"
					aritcle.Post("/", controllers.CreateArticle).Name = "创建文章"
					aritcle.Put("/{id:uint}", controllers.UpdateArticle).Name = "编辑文章"
					aritcle.Delete("/{id:uint}", controllers.DeleteArticle).Name = "删除文章"
				})

				admin.PartyFunc("/users", func(users iris.Party) {
					users.Get("/", controllers.GetAllUsers).Name = "用户列表"
					users.Get("/{id:uint}", controllers.GetUser).Name = "用户详情"
					users.Post("/", controllers.CreateUser).Name = "创建用户"
					users.Put("/{id:uint}", controllers.UpdateUser).Name = "编辑用户"

					users.Delete("/{id:uint}", controllers.DeleteUser).Name = "删除用户"
				})
				admin.PartyFunc("/roles", func(roles iris.Party) {
					roles.Get("/", controllers.GetAllRoles).Name = "角色列表"
					roles.Get("/{id:uint}", controllers.GetRole).Name = "角色详情"
					roles.Post("/", controllers.CreateRole).Name = "创建角色"
					roles.Put("/{id:uint}", controllers.UpdateRole).Name = "编辑角色"
					roles.Delete("/{id:uint}", controllers.DeleteRole).Name = "删除角色"
				})
				admin.PartyFunc("/permissions", func(permissions iris.Party) {
					permissions.Get("/", controllers.GetAllPermissions).Name = "权限列表"
					permissions.Get("/{id:uint}", controllers.GetPermission).Name = "权限详情"
					permissions.Post("/", controllers.CreatePermission).Name = "创建权限"
					permissions.Put("/{id:uint}", controllers.UpdatePermission).Name = "编辑权限"
					permissions.Delete("/{id:uint}", controllers.DeletePermission).Name = "删除权限"
				})
				admin.PartyFunc("/configs", func(configs iris.Party) {
					configs.Get("/", controllers.GetAllConfigs).Name = "系统配置列表"
					configs.Get("/{key:string}", controllers.GetConfig).Name = "系统配置详情"
					configs.Post("/", controllers.CreateConfig).Name = "创建系统配置"
					configs.Put("/{id:uint}", controllers.UpdateConfig).Name = "编辑系统配置"
					configs.Delete("/{id:uint}", controllers.DeleteConfig).Name = "删除系统配置"
				})
				admin.PartyFunc("/tags", func(tags iris.Party) {
					tags.Get("/", controllers.GetAllTags).Name = "标签列表"
					tags.Get("/{id:uint}", controllers.GetTag).Name = "标签详情"
					tags.Post("/", controllers.CreateTag).Name = "创建标签"
					tags.Put("/{id:uint}", controllers.UpdateTag).Name = "编辑标签"
					tags.Delete("/{id:uint}", controllers.DeleteTag).Name = "删除标签"
				})
				admin.PartyFunc("/types", func(types iris.Party) {
					types.Get("/", controllers.GetAllTypes).Name = "分类列表"
					types.Get("/{id:uint}", controllers.GetType).Name = "分类详情"
					types.Post("/", controllers.CreateType).Name = "创建分类"
					types.Put("/{id:uint}", controllers.UpdateType).Name = "编辑分类"
					types.Delete("/{id:uint}", controllers.DeleteType).Name = "删除分类"
				})
				admin.PartyFunc("/docs", func(docs iris.Party) {
					docs.Get("/", controllers.GetAllDocs).Name = "文档列表"
					docs.Get("/{id:uint}", controllers.GetDoc).Name = "文档详情"
					docs.Post("/", controllers.CreateDoc).Name = "创建文档"
					docs.Put("/{id:uint}", controllers.UpdateDoc).Name = "编辑文档"
					docs.Delete("/{id:uint}", controllers.DeleteDoc).Name = "删除文档"
				})
				admin.PartyFunc("/chapters", func(chapters iris.Party) {
					chapters.Get("/", controllers.GetAllChapters).Name = "章节列表"
					chapters.Get("/{id:uint}", controllers.GetChapter).Name = "章节详情"
					chapters.Post("/", controllers.CreateChapter).Name = "创建章节"
					chapters.Put("/{id:uint}", controllers.UpdateChapter).Name = "编辑章节"
					chapters.Put("/{id:uint}/set_sort", controllers.SetChapterSort).Name = "设置排序"
					chapters.Put("/sort", controllers.SortChapter).Name = "排序章节"
					chapters.Delete("/{id:uint}", controllers.DeleteChapter).Name = "删除章节"
				})
			})
		}
	}

}
