# 开发环境搭建指南

## 环境要求

- **Go**: 1.21 或更高版本
- **Node.js**: 16 或更高版本
- **npm**: 7 或更高版本
- **Git**: 用于版本控制

## 一、后端开发环境

### 1.1 安装Go

**macOS**
```bash
brew install go
```

**Linux (Ubuntu/Debian)**
```bash
sudo apt update
sudo apt install golang-go
```

**Windows**
下载安装包: https://go.dev/dl/

### 1.2 验证Go安装

```bash
go version
```

### 1.3 初始化后端项目

```bash
cd backend
go mod init poem/backend
```

### 1.4 安装依赖

```bash
# Gin框架
go get -u github.com/gin-gonic/gin

# CORS中间件
go get -u github.com/gin-contrib/cors

# Gzip压缩
go get -u github.com/gin-contrib/gzip
```

### 1.5 项目结构

创建以下目录结构：

```bash
mkdir -p api/handlers api/middleware
mkdir -p models services repository
mkdir -p config utils
```

### 1.6 运行开发服务器

```bash
cd backend
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 1.7 热重载开发

安装 `air` 实现热重载：

```bash
go install github.com/cosmtrek/air@latest
cd backend
air
```

## 二、前端开发环境

### 2.1 安装Node.js

**macOS**
```bash
brew install node
```

**Linux (Ubuntu/Debian)**
```bash
sudo apt update
sudo apt install nodejs npm
```

**Windows**
下载安装包: https://nodejs.org/

### 2.2 验证Node.js安装

```bash
node --version
npm --version
```

### 2.3 初始化Vue项目

```bash
cd frontend
npm create vite@latest . -- --template vue
```

### 2.4 安装依赖

```bash
# Vue Router
npm install vue-router@4

# Pinia状态管理
npm install pinia

# Axios HTTP客户端
npm install axios

# TailwindCSS
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

### 2.5 配置TailwindCSS

**tailwind.config.js**
```javascript
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'ink-black': '#1a1a1a',
        'rice-paper': '#f5f5f0',
        'cinnabar': '#c83c23',
      }
    },
  },
  plugins: [],
}
```

**src/assets/styles/main.css**
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

### 2.6 运行开发服务器

```bash
cd frontend
npm run dev
```

应用将在 `http://localhost:3000` 启动

## 三、开发工具推荐

### VS Code扩展

**后端开发**
- Go (golang.go)
- Error Lens (usernamehw.errorlens)

**前端开发**
- Volar (Vue.volar)
- ESLint (dbaeumer.vscode-eslint)
- Prettier (esbenp.prettier-vscode)

**通用**
- GitLens (eamodio.gitlens)
- Material Icon Theme (PKief.material-icon-theme)

### 配置文件

**.vscode/settings.json**
```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "golang.go"
  },
  "[vue]": {
    "editor.defaultFormatter": "Vue.volar"
  }
}
```

## 四、代码规范

### Go代码规范

遵循 [Effective Go](https://go.dev/doc/effective_go) 指南：

```go
// 包注释
// Package poetry 提供诗词相关的API服务
package poetry

// 导出函数注释
// GetPoems 获取诗词列表
// 参数:
//   - page: 页码
//   - pageSize: 每页数量
// 返回:
//   - poems: 诗词列表
//   - total: 总数
//   - err: 错误信息
func GetPoems(page, pageSize int) ([]Poem, int, error) {
    // 实现
}
```

### Vue代码规范

使用Composition API：

```vue
<template>
  <div class="poem-card">
    <h3>{{ poem.title }}</h3>
    <p>{{ poem.author }}</p>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  poem: {
    type: Object,
    required: true
  }
})
</script>

<style scoped>
.poem-card {
  padding: 1rem;
}
</style>
```

### 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| Go文件 | 小写+下划线 | poetry_service.go |
| Go结构体 | 大驼峰 | type PoetryService struct |
| Go变量 | 小驼峰 | poemCount |
| Vue文件 | 大驼峰 | PoemCard.vue |
| Vue组件 | 大驼峰 | <PoemCard /> |
| JS变量 | 小驼峰 | poemCount |

## 五、调试技巧

### 后端调试

使用 `fmt.Printf` 或 `log.Printf` 输出调试信息：

```go
log.Printf("获取诗词列表: page=%d, pageSize=%d", page, pageSize)
```

使用Delve调试器：

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go
```

### 前端调试

使用浏览器开发者工具：

```javascript
// Console输出
console.log('诗词数据:', poems)

// Vue Devtools
// 安装Vue Devtools浏览器扩展
```

### API测试

使用Postman或cURL测试API：

```bash
# 测试获取诗词列表
curl http://localhost:8080/api/v1/poems?page=1&page_size=10

# 测试获取随机诗词
curl http://localhost:8080/api/v1/poems/random
```

## 六、常见问题

### 问题1: Go模块依赖下载失败

**解决方案**: 设置Go代理

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 问题2: npm install速度慢

**解决方案**: 使用淘宝镜像

```bash
npm config set registry https://registry.npmmirror.com
```

### 问题3: CORS跨域问题

**解决方案**: 后端配置CORS中间件

```go
import "github.com/gin-contrib/cors"

router.Use(cors.Default())
```

### 问题4: 前端代理配置

在 `vite.config.js` 中配置代理：

```javascript
export default {
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
}
```

## 七、开发工作流

### 1. 创建新功能

```bash
# 创建功能分支
git checkout -b feature/new-feature

# 开发并测试
# ...

# 提交代码
git add .
git commit -m "feat: 添加新功能"
```

### 2. 代码审查清单

- [ ] 代码符合规范
- [ ] 添加必要注释
- [ ] 错误处理完善
- [ ] 无控制台错误
- [ ] 响应式布局正常

### 3. 测试流程

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm run test

# 集成测试
# 手动测试主要功能流程
```
