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
│   │   └── poetry.js           # 诗词状态
│   ├── views/                  # 页面组件
│   │   ├── HomeView.vue        # 首页
│   │   ├── CatalogView.vue     # 目录浏览
│   │   ├── AuthorsCatalogView.vue # 作者目录
│   │   ├── AuthorView.vue      # 作者详情页
│   │   ├── PoemDetailView.vue  # 诗词详情
│   │   └── SearchView.vue      # 搜索页面
│   ├── components/             # 可复用组件
│   │   ├── PoemCard.vue        # 诗词卡片
│   │   ├── PoemContent.vue     # 诗词内容
│   │   └── AuthorCard.vue      # 作者卡片
│   ├── composables/            # 组合式API
│   │   ├── usePoetry.js        # 诗词API逻辑
│   │   └── useSearch.js        # 搜索逻辑
│   ├── api/                    # API客户端
│   │   └── poetry-api.js
│   ├── utils/                  # 工具函数
│   │   └── common.js
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
| `/catalog` | CatalogView | 按分类浏览 |
| `/authors` | AuthorsCatalogView | 作者目录 |
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
    // ...
  }),
  actions: {
    // ...
  }
})
```

## 后端架构

### 目录结构

```
backend/
├── main.go                     # 应用入口
├── go.mod / go.sum             # 依赖管理
├── Dockerfile                  # 容器配置
├── config/                     # 配置
│   └── config.go
├── cmd/                        # 命令行工具
│   └── etl/
│       └── main.go             # ETL工具入口
├── api/                        # API层
│   ├── router.go               # 路由定义
│   ├── handlers/               # 请求处理器
│   │   └── poetry.go
│   └── middleware/             # 中间件
│       └── cors.go
├── models/                     # 数据模型
│   ├── poetry.go
│   └── poem_db.go
├── services/                   # 业务逻辑层
│   └── poetry_service.go
└── repository/                 # 数据访问层
    └── poetry_repository.go
```

### 分层架构

#### 1. API层 (handlers)
处理HTTP请求和响应，参数验证，调用Service层。
- `handlers/poetry.go`: 处理诗词相关的请求。

#### 2. Service层 (services)
业务逻辑处理，数据转换。
- `services/poetry_service.go`: 封装诗词业务逻辑。

#### 3. Repository层 (repository)
数据访问抽象，负责与数据库交互。
- `repository/poetry_repository.go`: 实现数据的增删改查。

#### 4. Model层 (models)
数据结构定义。
- `models/poetry.go`: 诗词相关的结构体定义。
- `models/poem_db.go`: 数据库相关的模型定义。

### 应用入口

```go
// main.go
package main

import (
    "poem/backend/api"
    "poem/backend/config"
    "poem/backend/repository"
    "poem/backend/services"
    // ...
)

func main() {
    // ...
}
```

## 数据流

### 诗词列表获取流程

```
用户请求 → 前端路由 → Vue组件 → usePoetry composable
    ↓
API客户端 (Axios)
    ↓
后端路由 → Handler (poetry.go) → Service (poetry_service.go) → Repository (poetry_repository.go)
    ↓
SQLite查询 (GORM)
    ↓
返回响应 → 前端状态更新 → 组件渲染
```

### 搜索流程

```
用户输入 → 搜索组件 → 防抖处理
    ↓
API搜索请求 → 后端SearchHandler
    ↓
SearchService → Repository (SQLite LIKE/FTS)
    ↓
返回结果 → 搜索结果展示
```

## 性能优化策略

### 前端优化
- 路由懒加载
- 组件复用
- 静态资源优化

### 后端优化
- 数据库索引优化
- 结构化日志
- 跨域资源共享 (CORS) 配置

## 安全考虑

- CORS配置
- 输入验证
- 错误处理与日志记录

## 扩展性设计

### 数据源扩展
Repository层抽象便于切换数据源，支持从JSON导入数据到SQLite。

### API扩展
模块化的API路由设计，便于添加新的功能模块（如用户系统、收藏功能等）。
