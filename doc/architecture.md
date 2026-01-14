# 系统架构设计

## 架构概览

本项目采用前后端分离架构，前端使用Vue.js，后端使用Go/Gin框架。

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Vue.js)                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Home      │  │  Catalog    │  │   Poem Detail       │ │
│  │   Page      │  │   Browser   │  │     Page            │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend (Go + Gin)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   API       │  │   Service   │  │    Data             │ │
│  │  Handlers   │  │    Layer    │  │    Access           │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Data Storage Layer                       │
│  ┌─────────────┐  ┌─────────────┐                           │
│  │    SQLite   │  │    JSON     │                           │
│  │   Database  │  │    Files    │                           │
│  └─────────────┘  └─────────────┘                           │
│        ▲               (ETL Source)                         │
│        │                                                    │
└────────┴────────────────────────────────────────────────────┘
```

## 前端架构

### 目录结构

```
frontend/
├── src/
│   ├── main.js                 # 应用入口
│   ├── App.vue                 # 根组件
│   ├── router/                 # 路由配置
│   │   └── index.js
│   ├── stores/                 # Pinia状态管理
│   │   ├── poetry.js           # 诗词状态
│   │   ├── author.js           # 作者状态
│   │   └── ui.js               # UI状态
│   ├── views/                  # 页面组件
│   │   ├── HomeView.vue        # 首页
│   │   ├── CatalogView.vue     # 目录浏览
│   │   ├── DynastyView.vue     # 朝代页面
│   │   ├── AuthorView.vue      # 作者页面
│   │   ├── PoemDetailView.vue  # 诗词详情
│   │   └── SearchView.vue      # 搜索页面
│   ├── components/             # 可复用组件
│   │   ├── PoemCard.vue        # 诗词卡片
│   │   ├── PoemContent.vue     # 诗词内容
│   │   ├── AuthorCard.vue      # 作者卡片
│   │   ├── SearchBar.vue       # 搜索栏
│   │   ├── Pagination.vue      # 分页组件
│   │   └── LoadingSpinner.vue  # 加载指示器
│   ├── composables/            # 组合式API
│   │   ├── usePoetry.js        # 诗词API
│   │   ├── useAuthor.js        # 作者API
│   │   └── useSearch.js        # 搜索API
│   ├── api/                    # API客户端
│   │   └── poetry-api.js
│   └── assets/                 # 静态资源
│       └── styles/
│           └── main.css
├── index.html
├── package.json
└── vite.config.js
```

### 路由设计

| 路径 | 组件 | 描述 |
|------|------|------|
| `/` | HomeView | 首页，展示精选诗词和快捷入口 |
| `/catalog` | CatalogView | 按朝代和分类浏览 |
| `/dynasty/:id` | DynastyView | 朝代详情页 |
| `/author/:id` | AuthorView | 作者详情页 |
| `/poem/:id` | PoemDetailView | 诗词详情页 |
| `/search` | SearchView | 搜索页面 |

### 状态管理

使用Pinia进行状态管理：

```javascript
// stores/poetry.js
import { defineStore } from 'pinia'

export const usePoetryStore = defineStore('poetry', {
  state: () => ({
    poems: [],
    currentPoem: null,
    totalCount: 0
  }),
  actions: {
    async fetchPoems(params) {
      // 获取诗词列表
    },
    async fetchPoemById(id) {
      // 获取单首诗词
    }
  }
})
```

## 后端架构

### 目录结构

```
backend/
├── main.go                     # 应用入口
├── go.mod / go.sum             # 依赖管理
├── config/                     # 配置
│   └── config.go
├── api/                        # API层
│   ├── router.go               # 路由定义
│   ├── handlers/               # 请求处理器
│   │   ├── poetry.go
│   │   ├── author.go
│   │   ├── dynasty.go
│   │   └── search.go
│   └── middleware/             # 中间件
│       ├── cors.go
│       └── logger.go
├── models/                     # 数据模型
│   ├── poetry.go
│   ├── author.go
│   ├── dynasty.go
│   └── response.go
├── services/                   # 业务逻辑层
│   ├── poetry_service.go
│   ├── author_service.go
│   ├── search_service.go
│   └── cache_service.go
└── repository/                 # 数据访问层
    ├── poetry_repository.go
    ├── author_repository.go
    └── json_loader.go
```

### 分层架构

#### 1. API层 (handlers)
处理HTTP请求和响应，参数验证，调用Service层

#### 2. Service层 (services)
业务逻辑处理，数据转换，缓存管理

#### 3. Repository层 (repository)
数据访问抽象，JSON文件读取，数据查询

#### 4. Model层 (models)
数据结构定义

### 应用入口

```go
// main.go
package main

import (
    "log"
    "poem/backend/api"
    "poem/backend/config"
    "poem/backend/repository"
    "poem/backend/services"
)

func main() {
    cfg := config.Load()

    // 初始化Repository层
    poetryRepo := repository.NewPoetryRepository(cfg.DataPath)
    authorRepo := repository.NewAuthorRepository(cfg.DataPath)

    // 初始化Service层
    poetryService := services.NewPoetryService(poetryRepo, authorRepo)
    searchService := services.NewSearchService(poetryRepo)

    // 初始化API层
    router := api.SetupRouter(poetryService, searchService)

    log.Printf("服务器启动在端口 %s", cfg.Port)
    router.Run(":" + cfg.Port)
}
```

## 数据流

### 诗词列表获取流程

```
用户请求 → 前端路由 → Vue组件 → usePoetry composable
    ↓
API客户端 (Axios)
    ↓
后端路由 → Handler → Service → Repository
    ↓
SQLite查询 (GORM)
    ↓
返回响应 → 前端状态更新 → 组件渲染
```

### 搜索流程

```
用户输入 → SearchBar组件 → 防抖处理
    ↓
API搜索请求 → 后端SearchHandler
    ↓
SearchService → Repository (SQLite LIKE/FTS)
    ↓
返回结果 → SearchResults组件
```

## 性能优化策略

### 前端优化
- 路由懒加载
- 虚拟滚动（长列表）
- 防抖搜索
- 组件缓存
- 代码分割

### 后端优化
- 数据库索引优化
- 列表查询分页
- GORM预加载 (Preload) 减少N+1查询
- 响应压缩（gzip）
- 连接池复用

## 安全考虑

- CORS配置
- SQL注入防护（GORM参数化查询）
- 输入验证
- HTTPS（生产环境）

## 扩展性设计

### 数据源扩展
Repository层抽象便于切换数据源：
- JSON文件（当前）
- SQLite数据库
- PostgreSQL/MySQL
- MongoDB

### 缓存扩展
CacheService接口设计支持：
- 内存缓存（当前）
- Redis
- Memcached

### API扩展
版本化API设计 (`/api/v1/`)，便于未来升级
