# 中国古诗词Web应用

一个展示中国古典诗词的Web应用，前端使用Vue.js，后端使用Go/Gin框架，数据来源于[chinese-poetry](https://github.com/chinese-poetry/chinese-poetry)数据集。

## 项目简介

本项目旨在构建一个优雅的中国古诗词展示平台，收录了：
- 全唐诗（约5.5万首）
- 全宋词（约21万首）
- 诗经、楚辞、元曲等古典文学作品
- 总计超过55万首诗词

## 技术栈

### 前端
- **Vue.js 3** - 渐进式JavaScript框架
- **Vite** - 新一代前端构建工具
- **Vue Router 4** - 官方路由管理器
- **Pinia** - 状态管理库
- **TailwindCSS** - 实用优先的CSS框架
- **Axios** - HTTP客户端

### 后端
- **Go 1.21+** - 高性能编程语言
- **Gin** - HTTP Web框架
- **CORS** - 跨域资源共享中间件

### 数据
- **chinese-poetry** - 中国古诗词JSON数据集

## 快速开始

### 环境要求

- Node.js 16+
- Go 1.21+
- Docker (可选)

### 1. 克隆项目

```bash
git clone <repository-url>
cd poem
```

### 2. 启动后端

```bash
cd backend
go mod download
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端应用将在 `http://localhost:3000` 启动

### 4. 使用Docker启动

```bash
docker-compose up -d
```

访问 `http://localhost` 查看应用

## 项目结构

```
poem/
├── chinese-poetry-master/    # 数据集（只读）
├── backend/                  # Go后端
│   ├── api/                  # API层
│   ├── models/               # 数据模型
│   ├── services/             # 业务逻辑
│   └── repository/           # 数据访问层
├── frontend/                 # Vue.js前端
│   └── src/
│       ├── views/            # 页面组件
│       ├── components/       # 可复用组件
│       ├── composables/      # 组合式API
│       └── api/              # API客户端
└── doc/                      # 项目文档
```

## 功能特性

- 诗词浏览与搜索
- 按朝代、作者分类
- 诗词详情展示
- 响应式设计
- 优雅的中国传统风格UI

## 文档

- [系统架构设计](architecture.md) - 完整的系统架构说明
- [API接口文档](api-reference.md) - RESTful API详细说明
- [开发指南](development-guide.md) - 开发环境搭建与规范
- [部署指南](deployment-guide.md) - 生产环境部署说明
- [数据模型说明](data-models.md) - 数据结构详解

## 许可证

MIT License

## 致谢

- [chinese-poetry](https://github.com/chinese-poetry/chinese-poetry) - 提供中国古诗词数据集
