package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/glebarez/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run scripts/migrate.go <数据库路径>")
		os.Exit(1)
	}

	dbPath := os.Args[1]

	// 连接数据库
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 尝试在不同路径查找SQL文件
	possiblePaths := []string{
		"migrations/create_users_table.sql",
		"../migrations/create_users_table.sql",
		"backend/migrations/create_users_table.sql",
	}

	var sqlBytes []byte
	var readErr error
	var foundPath string

	for _, p := range possiblePaths {
		sqlBytes, readErr = os.ReadFile(p)
		if readErr == nil {
			foundPath = p
			break
		}
	}

	if foundPath == "" {
		log.Fatalf("读取SQL文件失败: 无法找到 migrations/create_users_table.sql (尝试路径: %v)", possiblePaths)
	}

	fmt.Printf("正在使用迁移文件: %s\n", foundPath)

	// 执行SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("执行SQL失败:", err)
	}

	fmt.Println("✅ 用户表迁移成功!")
}
