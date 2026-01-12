package config

import (
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	Port     string
	DataPath string
	Env      string
}

// Load 加载配置
func Load() *Config {
	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		rootDir = "."
	}

	// 向上查找项目根目录
	for {
		if _, err := os.Stat(filepath.Join(rootDir, "chinese-poetry-master")); err == nil {
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			break
		}
		rootDir = parent
	}

	dataPath := filepath.Join(rootDir, "chinese-poetry-master")

	return &Config{
		Port:     getEnv("PORT", "8080"),
		DataPath: getEnv("DATA_PATH", dataPath),
		Env:      getEnv("ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
