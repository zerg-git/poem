package main

import (
	"log"
	"poem/backend/api"
	"poem/backend/config"
	"poem/backend/repository"
	"poem/backend/services"
)

func main() {
	// 加载配置
	cfg := config.Load()

	log.Printf("启动中国古诗词API服务...")
	log.Printf("数据路径: %s", cfg.DataPath)
	log.Printf("数据库路径: %s", cfg.DBPath)
	log.Printf("运行环境: %s", cfg.Env)

	// 初始化Repository层
	poetryRepo, db, err := repository.NewPoetryRepository(cfg.DBPath)
	if err != nil {
		log.Fatal("初始化数据库失败:", err)
	}

	// 初始化Service层
	poetryService := services.NewPoetryService(poetryRepo)

	// 设置路由
	router := api.SetupRouter(poetryService, db)

	// 启动服务器
	addr := ":" + cfg.Port
	log.Printf("服务器启动在 http://localhost%s", addr)
	log.Printf("API文档: http://localhost%s/api/v1", addr)
	log.Printf("API v2: http://localhost%s/api/v2", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
