# 中国古诗词Web应用

一个展示中国古典诗词的Web应用，前端使用Vue.js，后端使用Go/Gin框架。

![Vue.js](https://img.shields.io/badge/Vue.js-3.3-4FC08D?style=flat&logo=vue.js&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web-008ECF?style=flat)

## 项目简介

本项目基于 [chinese-poetry](https://github.com/chinese-poetry/chinese-poetry) 数据集，收录了超过55万首中国古典诗词，包括：
- 全唐诗（约5.5万首）
- 全宋词（约21万首）
- 诗经、楚辞、元曲等古典文学作品

## 功能特性

- 诗词浏览与分类导航
- 按朝代、作者分类展示
- 诗词详情页面
- 全文搜索功能
- 响应式设计
- 优雅的中国传统风格UI

## 快速开始

### 使用Docker（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd poem

# 启动服务
docker-compose up -d

# 访问应用
# 前端: http://localhost
# 后端API: http://localhost:8080/api/v1
```

### 手动启动

#### 后端

```bash
cd backend
go mod download
go run main.go
```

#### 前端

```bash
cd frontend
npm install
npm run dev
```

## 项目结构

```
poem/
├── chinese-poetry-master/    # 数据集（需要单独克隆）
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
├── doc/                      # 项目文档
└── docker-compose.yml        # Docker编排
```

## 准备数据集

由于数据集较大，需要单独下载：

```bash
cd poem
git clone https://github.com/chinese-poetry/chinese-poetry.git chinese-poetry-master
```

## API文档

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/api/v1/poems` | 获取诗词列表 |
| GET | `/api/v1/poems/:id` | 获取单首诗词 |
| GET | `/api/v1/poems/random` | 获取随机诗词 |
| GET | `/api/v1/authors` | 获取作者列表 |
| GET | `/api/v1/search` | 搜索诗词 |
| GET | `/api/v1/dynasties` | 获取朝代列表 |
| GET | `/api/v1/categories` | 获取分类列表 |

详细API文档请查看 [doc/api-reference.md](doc/api-reference.md)

## 技术栈

### 前端
- Vue.js 3 - 渐进式JavaScript框架
- Vue Router 4 - 路由管理
- Pinia - 状态管理
- Vite - 构建工具
- Axios - HTTP客户端

### 后端
- Go 1.21+ - 编程语言
- Gin - HTTP Web框架
- CORS - 跨域支持

## 开发

详细开发文档请查看 [doc/development-guide.md](doc/development-guide.md)

## 部署

详细部署文档请查看 [doc/deployment-guide.md](doc/deployment-guide.md)

## 许可证

MIT License

## 致谢

- [chinese-poetry](https://github.com/chinese-poetry/chinese-poetry) - 提供中国古诗词数据集
