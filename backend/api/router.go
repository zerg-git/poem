package api

import (
	"poem/backend/api/handlers"
	"poem/backend/api/middleware"
	"poem/backend/services"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(poetryService *services.PoetryService) *gin.Engine {
	router := gin.Default()

	// 使用中间件
	router.Use(middleware.CORS())

	// 创建处理器
	poetryHandler := handlers.NewPoetryHandler(poetryService)

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 目录相关
		v1.GET("/dynasties", poetryHandler.GetDynasties)
		v1.GET("/categories", poetryHandler.GetCategories)

		// 诗词相关
		v1.GET("/poems", poetryHandler.GetPoems)
		v1.GET("/poems/:id", poetryHandler.GetPoemByID)
		v1.GET("/poems/random", poetryHandler.GetRandomPoem)

		// 作者相关
		v1.GET("/authors", poetryHandler.GetAuthors)
		v1.GET("/authors/:name", poetryHandler.GetAuthorByName)
		v1.GET("/authors/:name/poems", poetryHandler.GetAuthorPoems)

		// 搜索
		v1.GET("/search", poetryHandler.Search)
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}
