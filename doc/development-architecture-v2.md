# 中国古诗词平台 - 开发架构设计文档 v2.0

## 文档版本信息
- **版本**: v2.0
- **创建日期**: 2026-01-15
- **文档状态**: 详细设计
- **基于项目**: d:\demo\poem

---

## 目录

1. [系统架构设计](#1-系统架构设计)
2. [数据库设计](#2-数据库设计)
3. [API接口设计](#3-api接口设计)
4. [前端架构设计](#4-前端架构设计)
5. [后端架构设计](#5-后端架构设计)
6. [核心功能实现方案](#6-核心功能实现方案)
7. [部署架构](#7-部署架构)
8. [开发计划](#8-开发计划)

---

## 1. 系统架构设计

### 1.1 整体架构概览

本项目采用**微服务化单体架构**，在保持部署简洁性的同时，通过清晰的模块边界支持未来扩展。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              客户端层                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   Web前端    │  │  移动端H5    │  │  卡片分享    │  │  管理后台    │   │
│  │ (Vue3+Vite)  │  │   (响应式)   │  │   (静态)     │  │  (未来扩展)  │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API网关层                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                      │
│  │  路由转发    │  │  鉴权中间件  │  │  限流中间件  │                      │
│  └──────────────┘  └──────────────┘  └──────────────┘                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                 ▼                 ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│   核心服务模块    │  │   AI服务模块     │  │   游戏服务模块   │
│  ┌──────────────┐ │  │  ┌──────────────┐ │  │  ┌──────────────┐ │
│  │ 诗词服务    │ │  │  │ AI解读服务  │ │  │  │ 飞花令对战   │ │
│  │ 作者服务    │ │  │  │ 智能问答    │ │  │  │ 每日挑战     │ │
│  │ 搜索服务    │ │  │  │ 学习路径    │ │  │  │ 成就系统     │ │
│  │ 推荐服务    │ │  │  └──────────────┘ │  │  │ 排行榜       │ │
│  │ 用户服务    │ │  │  ┌──────────────┐ │  │  └──────────────┘ │
│  │ 卡片服务    │ │  │  │ AI卡片生成  │ │  │  ┌──────────────┐ │
│  └──────────────┘ │  │  └──────────────┘ │  │  │ WebSocket   │ │
└──────────────────┘  └──────────────────┘  │  │  服务        │ │
                                            └──────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                 ▼                 ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│   SQLite数据库   │  │   Redis缓存      │  │   AI外部服务     │
│  ┌──────────────┐ │  │  ┌──────────────┐ │  │  ┌──────────────┐ │
│  │ 诗词数据    │ │  │  │ 会话存储    │ │  │  │ LLM API      │ │
│  │ 用户数据    │ │  │  │ 推荐缓存    │ │  │  │ 图像生成API  │ │
│  │ 学习记录    │ │  │  │ 游戏状态    │ │  │  │ (DALL-E/MJ)  │ │
│  │ 游戏数据    │ │  │  │ 实时排名    │ │  │  └──────────────┘ │
│  └──────────────┘ │  │  └──────────────┘ │  │  ┌──────────────┐ │
└──────────────────┘  └──────────────────┘  │  │ 向量数据库    │ │
                                            └──────────────────┘
```

### 1.2 技术栈选型

#### 前端技术栈
| 技术 | 版本 | 用途 |
|------|------|------|
| Vue.js | 3.3+ | 核心框架 |
| Vite | 4.4+ | 构建工具 |
| Vue Router | 4.2+ | 路由管理 |
| Pinia | 2.1+ | 状态管理 |
| Axios | 1.5+ | HTTP客户端 |
| **新增**: Socket.io-client | 2.5+ | 实时通信 |
| **新增**: html2canvas | 1.4+ | 卡片生成 |
| **新增**: Konva.js | 9.2+ | 图像处理 |
| **新增**: Chart.js | 4.4+ | 数据可视化 |
| **新增**: VueUse | 10.7+ | 组合式工具库 |

#### 后端技术栈
| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 核心语言 |
| Gin | 1.9+ | Web框架 |
| GORM | 1.31+ | ORM框架 |
| SQLite | 3.40+ | 主数据库 |
| **新增**: Redis | 7.2+ | 缓存与会话 |
| **新增**: WebSocket | - | 实时通信 |
| **新增**: go-socket.io | 1.7+ | Socket.IO服务端 |
| **新增**: jwt-go | 5.2+ | JWT认证 |
| **新增**: go-openai | - | AI服务集成 |
| **新增**: gorse | - | 推荐系统框架 |

#### AI服务集成
| 服务 | 用途 |
|------|------|
| OpenAI GPT-4 API | 诗词解读、智能问答 |
| 智谱AI GLM-4 API | 中文优化解读 |
| DALL-E 3 / Midjourney | 卡片背景图生成 |
| 向量数据库 | 诗词语义检索 |

### 1.3 服务拆分方案

#### 阶段一：模块化单体（当前推荐）
保持单一部署单元，通过包级别隔离实现模块化：
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── core/          # 核心诗词服务
│   ├── ai/            # AI服务模块
│   ├── game/          # 游戏模块
│   ├── user/          # 用户模块
│   ├── recommend/     # 推荐模块
│   └── card/          # 卡片生成模块
├── pkg/               # 共享包
└── api/               # API层
```

---

## 2. 数据库设计

### 2.1 新增数据表设计

#### 用户相关表

```sql
-- 用户表
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    open_id VARCHAR(100) UNIQUE,
    union_id VARCHAR(100),
    username VARCHAR(50) UNIQUE NOT NULL,
    nickname VARCHAR(100),
    avatar_url VARCHAR(500),
    email VARCHAR(100),
    phone VARCHAR(20),
    gender TINYINT DEFAULT 0,
    birth_date DATE,
    province VARCHAR(50),
    city VARCHAR(50),
    level INTEGER DEFAULT 1,
    experience INTEGER DEFAULT 0,
    coins INTEGER DEFAULT 0,
    vip_level INTEGER DEFAULT 0,
    vip_expire_at DATETIME,
    status TINYINT DEFAULT 1,
    last_login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_open_id ON users(open_id);
CREATE INDEX idx_users_username ON users(username);
```

#### 学习系统表

```sql
-- 学习路径表
CREATE TABLE learning_paths (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    difficulty TINYINT DEFAULT 1,
    dynasty VARCHAR(50),
    category VARCHAR(50),
    estimated_hours INTEGER,
    poem_count INTEGER,
    cover_url VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 学习路径诗词关联表
CREATE TABLE path_poems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path_id INTEGER NOT NULL,
    work_id INTEGER NOT NULL,
    sort_order INTEGER DEFAULT 0,
    FOREIGN KEY (path_id) REFERENCES learning_paths(id),
    FOREIGN KEY (work_id) REFERENCES works(id),
    UNIQUE(path_id, work_id)
);

-- 用户学习记录表
CREATE TABLE user_learning_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    work_id INTEGER NOT NULL,
    path_id INTEGER,
    status TINYINT DEFAULT 0,
    reading_count INTEGER DEFAULT 0,
    favorite_count INTEGER DEFAULT 0,
    last_read_at DATETIME,
    mastery_level INTEGER DEFAULT 0,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (work_id) REFERENCES works(id),
    FOREIGN KEY (path_id) REFERENCES learning_paths(id),
    UNIQUE(user_id, work_id)
);

-- AI解读缓存表
CREATE TABLE ai_interpretations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    work_id INTEGER NOT NULL UNIQUE,
    translation TEXT,
    annotation TEXT,
    analysis TEXT,
    difficulty_score DECIMAL(3,2),
    tags VARCHAR(500),
    model_version VARCHAR(50),
    cached_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (work_id) REFERENCES works(id)
);

-- 智能问答历史表
CREATE TABLE qa_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    session_id VARCHAR(100) NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    related_poems TEXT,
    feedback TINYINT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_qa_session ON qa_history(session_id);
```

#### 游戏系统表

```sql
-- 飞花令游戏表
CREATE TABLE feihualing_games (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id VARCHAR(50) UNIQUE NOT NULL,
    game_type TINYINT DEFAULT 1,
    keyword VARCHAR(10) NOT NULL,
    difficulty TINYINT DEFAULT 1,
    status TINYINT DEFAULT 0,
    current_turn INTEGER DEFAULT 0,
    max_turns INTEGER DEFAULT 10,
    created_by INTEGER,
    winner_id INTEGER,
    started_at DATETIME,
    ended_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (winner_id) REFERENCES users(id)
);

-- 飞花令参与记录表
CREATE TABLE feihualing_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TINYINT DEFAULT 1,
    score INTEGER DEFAULT 0,
    correct_count INTEGER DEFAULT 0,
    wrong_count INTEGER DEFAULT 0,
    total_time INTEGER DEFAULT 0,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES feihualing_games(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(game_id, user_id)
);

-- 飞花令出招记录表
CREATE TABLE feihualing_moves (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    turn_number INTEGER NOT NULL,
    work_id INTEGER NOT NULL,
    poem_line TEXT NOT NULL,
    response_time INTEGER,
    is_valid TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES feihualing_games(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (work_id) REFERENCES works(id)
);

-- 每日挑战表
CREATE TABLE daily_challenges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    challenge_date DATE UNIQUE NOT NULL,
    challenge_type TINYINT DEFAULT 1,
    content TEXT NOT NULL,
    difficulty TINYINT DEFAULT 1,
    hints TEXT,
    answers TEXT NOT NULL,
    completion_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 成就表
CREATE TABLE achievements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    category VARCHAR(50),
    rarity TINYINT DEFAULT 1,
    requirement_type VARCHAR(50),
    requirement_value INTEGER,
    reward_coins INTEGER DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户成就表
CREATE TABLE user_achievements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    achievement_id INTEGER NOT NULL,
    progress INTEGER DEFAULT 0,
    unlocked_at DATETIME,
    notified TINYINT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (achievement_id) REFERENCES achievements(id),
    UNIQUE(user_id, achievement_id)
);
```

#### 推荐系统表

```sql
-- 用户行为记录表
CREATE TABLE user_behaviors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    behavior_type TINYINT NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_id INTEGER NOT NULL,
    duration INTEGER,
    context TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_behaviors_user ON user_behaviors(user_id);
CREATE INDEX idx_behaviors_type ON user_behaviors(behavior_type);
CREATE INDEX idx_behaviors_created ON user_behaviors(created_at);

-- 推荐缓存表
CREATE TABLE recommendations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    recommend_type TINYINT NOT NULL,
    content TEXT NOT NULL,
    generated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expire_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### 卡片系统表

```sql
-- 卡片模板表
CREATE TABLE card_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    preview_url VARCHAR(500),
    template_config TEXT NOT NULL,
    category VARCHAR(50),
    is_premium TINYINT DEFAULT 0,
    usage_count INTEGER DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户生成的卡片表
CREATE TABLE user_cards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    template_id INTEGER NOT NULL,
    work_id INTEGER NOT NULL,
    card_config TEXT NOT NULL,
    image_url VARCHAR(500),
    background_prompt TEXT,
    is_public TINYINT DEFAULT 0,
    share_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (template_id) REFERENCES card_templates(id),
    FOREIGN KEY (work_id) REFERENCES works(id)
);

CREATE INDEX idx_cards_user ON user_cards(user_id);
CREATE INDEX idx_cards_public ON user_cards(is_public);
```

### 2.2 索引优化策略

```sql
-- 诗词查询优化
CREATE INDEX idx_works_title ON works(title);
CREATE INDEX idx_works_author ON works(author_id);
CREATE INDEX idx_works_category ON works(category_id);

-- 全文搜索优化（SQLite FTS5）
CREATE VIRTUAL TABLE works_fts USING fts5(
    title, content, prologue,
    content='works',
    content_rowid='id'
);

-- 用户行为分析优化
CREATE INDEX idx_behaviors_composite ON user_behaviors(user_id, behavior_type, created_at);
```

---

## 3. API接口设计

### 3.1 RESTful API设计规范

#### 通用响应格式
```json
{
  "code": 200,
  "success": true,
  "data": {},
  "message": "操作成功",
  "timestamp": 1704067200000
}
```

### 3.2 各模块API端点定义

#### AI诗词助手模块
```
GET    /api/v2/ai/poems/:id/interpretation      # 获取诗词AI解读
POST   /api/v2/ai/poems/:id/interpretation      # 重新生成解读
GET    /api/v2/ai/poems/:id/translation         # 获取白话翻译
POST   /api/v2/ai/qa/chat                       # 发起问答
GET    /api/v2/ai/qa/history                    # 问答历史
GET    /api/v2/learning/paths                   # 学习路径列表
GET    /api/v2/learning/paths/:id               # 路径详情
POST   /api/v2/learning/progress                # 更新学习进度
GET    /api/v2/learning/dashboard               # 学习数据看板
```

#### 诗词飞花令模块
```
POST   /api/v2/game/feihualing/rooms            # 创建游戏房间
GET    /api/v2/game/feihualing/rooms/:id        # 房间详情
POST   /api/v2/game/feihualing/rooms/:id/join   # 加入房间
POST   /api/v2/game/feihualing/rooms/:id/move   # 提交诗句
GET    /api/v2/game/feihualing/ranking          # 排行榜
GET    /api/v2/game/daily/today                 # 今日挑战
GET    /api/v2/game/achievements                # 成就列表
GET    /api/v2/game/achievements/mine           # 我的成就
```

#### 卡片生成模块
```
GET    /api/v2/cards/templates                  # 模板列表
POST   /api/v2/cards/generate                   # 生成卡片
POST   /api/v2/cards/generate/async             # 异步生成
GET    /api/v2/cards/:id                        # 卡片详情
GET    /api/v2/cards/mine                       # 我的卡片
GET    /api/v2/cards/gallery                    # 卡片广场
POST   /api/v2/cards/:id/like                   # 点赞卡片
```

#### 推荐系统模块
```
GET    /api/v2/recommend/poems                  # 推荐诗词
GET    /api/v2/recommend/authors                # 推荐作者
GET    /api/v2/recommend/discover               # 发现页
POST   /api/v2/recommend/feedback               # 反馈推荐结果
```

#### 用户系统模块
```
POST   /api/v2/auth/register                    # 注册
POST   /api/v2/auth/login                       # 登录
POST   /api/v2/auth/logout                      # 登出
GET    /api/v2/users/profile                    # 个人资料
PUT    /api/v2/users/profile                    # 更新资料
POST   /api/v2/users/favorites                  # 收藏诗词
GET    /api/v2/users/favorites                  # 收藏列表
DELETE /api/v2/users/favorites/:id              # 取消收藏
```

### 3.3 WebSocket接口设计

```
WS     /api/v2/ws/feihualing/:room_id

# 消息类型
{
  "type": "move|join|leave|game_start|game_end|timer",
  "data": {},
  "timestamp": 1704067200000
}
```

---

## 4. 前端架构设计

### 4.1 目录结构

```
frontend/
├── src/
│   ├── views/
│   │   ├── core/                    # 核心页面
│   │   ├── ai/                      # AI模块
│   │   │   ├── AIInterpretView.vue
│   │   │   ├── AIQAView.vue
│   │   │   ├── LearningPathsView.vue
│   │   │   └── LearningDashboard.vue
│   │   ├── game/                    # 游戏模块
│   │   │   ├── FeihualingLobby.vue
│   │   │   ├── FeihualingRoom.vue
│   │   │   ├── DailyChallengeView.vue
│   │   │   └── AchievementsView.vue
│   │   ├── card/                    # 卡片模块
│   │   │   ├── CardTemplatesView.vue
│   │   │   ├── CardGeneratorView.vue
│   │   │   └── CardGalleryView.vue
│   │   └── user/                    # 用户模块
│   ├── components/
│   │   ├── ai/                      # AI组件
│   │   ├── game/                    # 游戏组件
│   │   └── card/                    # 卡片组件
│   ├── stores/
│   │   ├── ai.js                    # AI状态
│   │   ├── game.js                  # 游戏状态
│   │   ├── card.js                  # 卡片状态
│   │   └── user.js                  # 用户状态
│   ├── composables/
│   │   ├── useAI.js
│   │   ├── useGame.js
│   │   ├── useCard.js
│   │   └── useWebSocket.js
│   └── api/
│       ├── ai-api.js
│       ├── game-api.js
│       └── card-api.js
```

### 4.2 状态管理

```javascript
// stores/ai.js - AI服务状态
export const useAIStore = defineStore('ai', () => {
  const interpretations = ref(new Map())
  const chatHistory = ref([])
  const learningProgress = ref(new Map())

  async function getInterpretation(poemId) {
    if (interpretations.value.has(poemId)) {
      return interpretations.value.get(poemId)
    }
    // API调用...
  }

  return { interpretations, chatHistory, learningProgress, getInterpretation }
})
```

### 4.3 路由设计

```javascript
// router/index.js
const routes = [
  { path: '/', component: HomeView },
  {
    path: '/ai',
    children: [
      { path: 'interpret/:id', component: AIInterpretView },
      { path: 'qa', component: AIQAView },
      { path: 'learn', component: LearningPathsView },
      { path: 'dashboard', component: LearningDashboard, meta: { requiresAuth: true } }
    ]
  },
  {
    path: '/game',
    children: [
      { path: 'feihualing', component: FeihualingLobby },
      { path: 'feihualing/:roomId', component: FeihualingRoom },
      { path: 'daily', component: DailyChallengeView },
      { path: 'achievements', component: AchievementsView }
    ]
  },
  {
    path: '/cards',
    children: [
      { path: 'generate', component: CardGeneratorView },
      { path: 'gallery', component: CardGalleryView }
    ]
  }
]
```

---

## 5. 后端架构设计

### 5.1 项目结构

```
backend/
├── internal/
│   ├── core/                    # 核心诗词服务
│   ├── ai/                      # AI服务模块
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   ├── game/                    # 游戏模块
│   │   ├── handler/
│   │   ├── service/
│   │   ├── websocket/
│   │   └── repository/
│   ├── user/                    # 用户模块
│   ├── recommend/               # 推荐模块
│   └── card/                    # 卡片模块
├── pkg/
│   ├── middleware/              # 中间件
│   ├── models/                  # 数据模型
│   └── cache/                   # 缓存
└── api/
    └── v2/                      # API v2
```

### 5.2 分层架构示例

```go
// Handler层
type InterpretationHandler struct {
    service *service.InterpretationService
}

func (h *InterpretationHandler) GetInterpretation(c *gin.Context) {
    poemID := c.Param("id")
    interpretation, err := h.service.GetInterpretation(c.Request.Context(), poemID)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.Success(c, interpretation)
}

// Service层
type InterpretationService struct {
    repo     repository.InterpretationRepository
    cache    cache.Cache
    aiClient *AIClient
}

func (s *InterpretationService) GetInterpretation(ctx context.Context, poemID string) (*models.AIInterpretation, error) {
    // 缓存检查
    if cached, err := s.cache.Get(ctx, "interpretation:"+poemId); err == nil {
        return cached, nil
    }
    // 数据库查询
    interpretation, err := s.repo.GetByPoemID(ctx, poemID)
    if err == nil {
        return interpretation, nil
    }
    // AI生成
    return s.aiClient.GenerateInterpretation(ctx, poemID)
}
```

### 5.3 中间件设计

```go
// JWT认证中间件
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
            c.Abort()
            return
        }
        // Token验证逻辑...
        c.Next()
    }
}

// 限流中间件
func RateLimit(rps int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(rps), rps*2)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## 6. 核心功能实现方案

### 6.1 AI集成方案

#### AI服务客户端
```go
type AIClient interface {
    GenerateInterpretation(ctx context.Context, poemID string) (*models.AIInterpretation, error)
    AnswerQuestion(ctx context.Context, question string, history []ChatMessage) (string, error)
    GenerateBackgroundImage(ctx context.Context, prompt string) (string, error)
}

// 多级缓存策略
// L1: 内存缓存（热门诗词1000条）
// L2: Redis缓存（24小时）
// L3: SQLite持久化
// L4: AI生成
```

#### Prompt设计
```
你是一位中国古诗词专家。请对以下诗词进行专业解读：

诗词标题：%s
作者：%s
朝代：%s
内容：%s

请提供：
1. 白话文翻译：通俗易懂的现代文翻译
2. 典故注释：解释诗中使用的典故和难点词汇
3. 赏析分析：分析诗词的艺术特色、思想内涵
4. 难度评分：0-1之间的数值
5. 标签：3-5个描述诗词特点的标签
```

### 6.2 推荐算法设计

#### 协同过滤实现
```go
type CollaborativeFilter struct {
    behaviorRepo repository.BehaviorRepository
}

// UserBasedCF 基于用户的协同过滤
func (cf *CollaborativeFilter) UserBasedCF(ctx context.Context, userID uint, n int) ([]int, error) {
    // 1. 计算用户相似度（余弦相似度）
    similarUsers := cf.findSimilarUsers(ctx, userID, 20)
    // 2. 基于相似用户的行为计算推荐分数
    scores := cf.calculateScores(similarUsers)
    // 3. 返回TopN
    return topN(scores, n), nil
}

// 混合推荐策略
// 60% 协同过滤 + 30% 内容推荐 + 10% 热门趋势
```

### 6.3 实时对战实现

#### WebSocket服务
```go
type FeihualingHub struct {
    rooms   map[string]*GameRoom
    clients map[*Client]bool
}

type GameRoom struct {
    ID          string
    Keyword     string
    Status      string
    Players     map[uint]*Player
    Moves       []Move
    CurrentTurn uint
}

func (h *FeihualingHub) HandleMove(client *Client, moveData map[string]interface{}) {
    // 1. 验证出招（包含关键字、不重复）
    // 2. 更新游戏状态
    // 3. 广播给所有玩家
    // 4. 检查游戏结束条件
}
```

### 6.4 卡片生成实现

#### 前端生成方案
```javascript
import html2canvas from 'html2canvas'

export class CardGenerator {
  async generateWithCanvas(poem, template, options) {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    // 1. 绘制背景
    await this.drawBackground(ctx, template.background)
    // 2. 绘制诗词内容
    this.drawPoemContent(ctx, poem, template.layout)
    // 3. 添加印章和装饰
    this.drawDecorations(ctx, template.decorations)
    return canvas.toDataURL('image/png')
  }
}
```

---

## 7. 部署架构

### 7.1 Docker容器化

#### docker-compose.yml
```yaml
version: '3.8'

services:
  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    networks:
      - poetry-network

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_PATH=/data/poems.db
      - REDIS_ADDR=redis:6379
    volumes:
      - poetry-data:/data
    depends_on:
      - redis
    networks:
      - poetry-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    networks:
      - poetry-network

networks:
  poetry-network:
    driver: bridge

volumes:
  poetry-data:
```

### 7.2 Nginx配置

```nginx
upstream backend {
    server backend:8080;
}

upstream frontend {
    server frontend:80;
}

server {
    listen 80;
    server_name poetry.example.com;

    # 前端
    location / {
        proxy_pass http://frontend;
    }

    # API
    location /api/ {
        proxy_pass http://backend;
        limit_req zone=api_limit burst=20 nodelay;
    }

    # WebSocket
    location /ws/ {
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## 8. 开发计划

### 8.1 迭代规划

| 阶段 | 任务 | 工作量 | 里程碑 |
|------|------|--------|--------|
| 第一阶段 | 用户系统、基础框架 | 4周 | M1: 基础版 |
| 第二阶段 | AI诗词助手 | 6周 | M2: AI助手 |
| 第三阶段 | 飞花令游戏 | 5周 | M3: 游戏上线 |
| 第四阶段 | 卡片生成系统 | 4周 | M4: 卡片系统 |
| 第五阶段 | 推荐系统 | 4周 | M5: 推荐系统 |
| 第六阶段 | 全功能集成 | 3周 | M6: 完整版 |

**总开发周期**: 约26周（6个月）

### 8.2 团队配置

| 角色 | 人数 | 主要职责 |
|------|------|----------|
| 前端工程师 | 2-3人 | Vue组件开发、UI实现 |
| 后端工程师 | 2-3人 | API开发、业务逻辑 |
| AI工程师 | 1人 | AI集成、算法调优 |
| UI/UX设计师 | 1人 | 界面设计、交互设计 |
| 测试工程师 | 1人 | 功能测试、自动化测试 |
| 产品经理 | 1人 | 需求管理、迭代规划 |

### 8.3 关键指标

| 指标 | 目标值 |
|------|--------|
| DAU增长 | +300% |
| 用户停留时长 | 3分钟 → 15分钟 |
| AI解读准确率 | >80% |
| 推荐点击率 | >15% |
| 日生成卡片数 | >500 |

---

## 附录

### A. 技术选型对比

#### AI服务供应商
| 供应商 | 优势 | 劣势 | 推荐用途 |
|--------|------|------|----------|
| OpenAI GPT-4 | 能力强、生态好 | 中文一般 | 高质量解读 |
| 智谱AI GLM-4 | 中文优化、价格合理 | 生态较小 | 主要使用 |
| 文心一言 | 国内稳定 | 能力稍弱 | 备用方案 |

### B. 风险评估

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| AI API限流 | 高 | 多供应商切换、缓存策略 |
| WebSocket不稳定 | 中 | 自动重连、心跳机制 |
| SQLite性能瓶颈 | 中 | 准备PostgreSQL迁移方案 |

---

## 文档总结

本架构设计文档基于现有的Vue3+Go技术栈，详细规划了四大核心功能的实现方案：

1. **AI诗词助手**: 集成OpenAI/智谱AI，实现智能解读、问答和学习路径
2. **飞花令游戏**: 基于WebSocket的实时对战系统
3. **卡片生成**: 结合AI图像生成的创意分享功能
4. **推荐系统**: 协同过滤+内容推荐的混合策略

**架构原则**:
- 渐进式增强：在现有架构基础上扩展
- 模块化设计：清晰的模块边界
- 性能优先：多级缓存、异步处理
- 可扩展性：支持水平扩展

**预计开发周期**: 26周（约6个月）
