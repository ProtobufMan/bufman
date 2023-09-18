package router

import (
	"github.com/ProtobufMan/bufman/internal/handlers/http_handlers"
	"github.com/gin-gonic/gin"
)

func InitHTTPRouter() *gin.Engine {
	router := gin.Default()

	authn := router.Group("/authn")
	{
		authn.GET("/current_user", http_handlers.AuthnGroup.GetCurrentUser) // 根据token获取当前用户信息
	}

	user := router.Group("/user")
	{
		user.POST("/create", http_handlers.UserGroup.CreateUser) // 创建用户
		user.GET("/:id", http_handlers.UserGroup.GetUser)        // 查询用户
		user.POST("/list", http_handlers.UserGroup.ListUsers)    // 批量查询用户
	}

	token := router.Group("/token")
	{
		token.POST("/create", http_handlers.TokenGroup.CreateToken)      // 创建token
		token.GET("/:token_id", http_handlers.TokenGroup.GetToken)       // 获取token
		token.POST("/list", http_handlers.TokenGroup.ListTokens)         // 批量查询token
		token.DELETE("/:token_id", http_handlers.TokenGroup.DeleteToken) // 删除tokens
	}

	repository := router.Group("/repository")
	{
		repository.POST("/create", http_handlers.RepositoryGroup.CreateRepositoryByFullName)                         // 创建repository
		repository.GET("/:id", http_handlers.RepositoryGroup.GetRepository)                                          // 根据id获取repository
		repository.GET("/:repository_owner/:repository_name", http_handlers.RepositoryGroup.GetRepositoryByFullName) // 获取repository
		repository.POST("/list", http_handlers.RepositoryGroup.ListRepositories)                                     // 批量查询所有repository
		repository.DELETE("/:id", http_handlers.RepositoryGroup.DeleteRepository)                                    // 删除repository
		repository.POST("/list/:user_id", http_handlers.RepositoryGroup.ListUserRepositories)                        // 批量查询用户的repository
		repository.POST("/list_accessible", http_handlers.RepositoryGroup.ListRepositoriesUserCanAccess)             // 批量查询当前用户可访问的repository
		repository.PUT("/deprecate", http_handlers.RepositoryGroup.DeprecateRepositoryByName)                        // 弃用repository
		repository.PUT("/undeprecate", http_handlers.RepositoryGroup.UndeprecateRepositoryByName)                    // 解除弃用
		repository.PUT("/update", http_handlers.RepositoryGroup.UpdateRepositorySettingsByName)                      // 更新repository

		commit := repository.Group("/commit")
		{
			commit.POST("/list/:repository_owner/:repository_name/:reference", http_handlers.CommitGroup.ListRepositoryCommitsByReference) // 获取reference对应commit以及之前的commits
			commit.GET("/:repository_owner/:repository_name/:reference", http_handlers.CommitGroup.GetRepositoryCommitByReference)         // 获取reference对应commit
			commit.POST("/draft/list/:repository_owner/:repository_name", http_handlers.CommitGroup.ListRepositoryDraftCommits)            // 获取所有的草稿
			commit.DELETE("/draft/:repository_owner/:repository_name/:draft_name", http_handlers.CommitGroup.DeleteRepositoryDraftCommit)  // 删除草稿
		}

		tag := repository.Group("/tag")
		{
			tag.POST("/create", http_handlers.TagGroup.CreateRepositoryTag) // 创建tag
			tag.POST("/list", http_handlers.TagGroup.ListRepositoryTags)    // 查询repository下的所有tag
		}

		doc := repository.Group("/doc")
		{
			doc.GET("/source/:repository_owner/:repository_name/:reference", http_handlers.DocGroup.GetSourceDirectoryInfo)                 // 获取目录信息
			doc.GET("/source/:repository_owner/:repository_name/:reference/:path", http_handlers.DocGroup.GetSourceFile)                    // 获取文件源码
			doc.GET("/module/:repository_owner/:repository_name/:reference", http_handlers.DocGroup.GetModuleDocumentation)                 // 获取repo说明文档
			doc.GET("/package/:repository_owner/:repository_name/:reference", http_handlers.DocGroup.GetModulePackages)                     // 获取repo packages
			doc.GET("/package/:repository_owner/:repository_name/:reference/:package_name", http_handlers.DocGroup.GetPackageDocumentation) //获取包说明文档
		}
	}

	plugin := router.Group("/plugin")
	{
		plugin.POST("/create", http_handlers.PluginGroup.CreateCuratedPlugin)                            // 创建插件
		plugin.POST("/list", http_handlers.PluginGroup.ListCuratedPlugins)                               // 批量查询插件
		plugin.GET("/:owner/:name/:version/:revision", http_handlers.PluginGroup.GetLatestCuratedPlugin) // 查询插件
	}

	return router
}
