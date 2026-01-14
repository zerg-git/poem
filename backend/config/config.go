package config

import (
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	Port     string
	DataPath string
	DBPath   string
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
	var dataDirName = "chinese-poetry-master"
	for {
		if _, err := os.Stat(filepath.Join(rootDir, "chinese-poetry-master")); err == nil {
			dataDirName = "chinese-poetry-master"
			break
		}
		if _, err := os.Stat(filepath.Join(rootDir, "chinese-poetry")); err == nil {
			dataDirName = "chinese-poetry"
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			break
		}
		rootDir = parent
	}

	dataPath := filepath.Join(rootDir, dataDirName)
	// 数据库路径默认为项目根目录下的 poems.db
	dbPath := filepath.Join(rootDir, "poems.db")

	return &Config{
		Port:     getEnv("PORT", "8080"),
		DataPath: getEnv("DATA_PATH", dataPath),
		DBPath:   getEnv("DB_PATH", dbPath),
		Env:      getEnv("ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
