# 用户登录注册功能实现文档

## 概述

本文档描述了中国古诗词平台用户登录注册功能的完整实现，包括后端API、前端页面和数据库设计。

## 技术栈

### 后端
- **语言**: Go 1.24+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt

### 前端
- **框架**: Vue 3
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **HTTP客户端**: Axios

## 数据库设计

### 用户表 (users)

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    open_id VARCHAR(100) UNIQUE,           -- 微信OpenID（预留）
    union_id VARCHAR(100),                 -- 微信UnionID（预留）
    username VARCHAR(50) NOT NULL UNIQUE,  -- 用户名
    password_hash VARCHAR(255) NOT NULL,    -- 密码哈希
    nickname VARCHAR(100),                 -- 昵称
    avatar_url VARCHAR(500),               -- 头像URL
    email VARCHAR(100) UNIQUE,             -- 邮箱
    phone VARCHAR(20) UNIQUE,              -- 手机号
    gender INTEGER DEFAULT 0,              -- 性别: 0未知 1男 2女
    birth_date DATE,                       -- 生日
    province VARCHAR(50),                  -- 省份
    city VARCHAR(50),                      -- 城市
    level INTEGER DEFAULT 1,               -- 等级
    experience INTEGER DEFAULT 0,          -- 经验值
    coins INTEGER DEFAULT 0,               -- 金币
    vip_level INTEGER DEFAULT 0,           -- VIP等级
    vip_expire_at DATETIME,                -- VIP过期时间
    status INTEGER DEFAULT 1,              -- 状态: 0禁用 1正常
    last_login_at DATETIME,                -- 最后登录时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 用户收藏表 (user_favorites)

```sql
CREATE TABLE user_favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,            -- 诗词ID或作者ID
    target_type VARCHAR(20) NOT NULL,      -- poem / author
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, target_id, target_type)
);
```

### 用户浏览历史表 (user_history)

```sql
CREATE TABLE user_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,            -- 诗词ID或作者ID
    target_type VARCHAR(20) NOT NULL,      -- poem / author
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## 后端实现

### 文件结构

```
backend/
├── models/
│   └── user.go                           # 用户数据模型
├── repository/
│   └── user_repository.go                # 用户数据访问层
├── services/user/
│   └── user_service.go                   # 用户业务逻辑层
├── pkg/
│   ├── auth/
│   │   ├── jwt.go                        # JWT管理器
│   │   └── password.go                   # 密码加密工具
│   └── response/
│       └── response.go                   # 统一响应格式
├── api/
│   ├── handlers/v2/
│   │   └── user.go                       # 用户API处理器
│   ├── middleware/
│   │   └── auth.go                       # JWT认证中间件
│   └── v2/
│       └── router.go                     # API v2路由
└── migrations/
    └── create_users_table.sql            # 数据库迁移脚本
```

### API接口

#### 1. 用户注册
```
POST /api/v2/auth/register

请求体:
{
  "username": "testuser",
  "password": "password123",
  "nickname": "测试用户",
  "email": "test@example.com"
}

响应:
{
  "code": 200,
  "success": true,
  "data": {
    "token": "eyJhbGc...",
    "expires_at": 1705334400,
    "user": {
      "id": 1,
      "username": "testuser",
      "nickname": "测试用户",
      ...
    }
  }
}
```

#### 2. 用户登录
```
POST /api/v2/auth/login

请求体:
{
  "username": "testuser",
  "password": "password123"
}

响应: (同注册)
```

#### 3. 获取个人资料
```
GET /api/v2/users/profile
Headers: Authorization: Bearer {token}

响应:
{
  "code": 200,
  "success": true,
  "data": {
    "id": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "level": 1,
    "experience": 0,
    "coins": 0,
    ...
  }
}
```

#### 4. 刷新Token
```
POST /api/v2/auth/refresh

请求体:
{
  "token": "eyJhbGc..."
}

响应:
{
  "code": 200,
  "success": true,
  "data": {
    "token": "eyJhbGc..."
  }
}
```

### 核心代码示例

#### JWT认证中间件
```go
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Unauthorized(c, "未登录")
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            response.Unauthorized(c, "Token格式错误")
            c.Abort()
            return
        }

        claims, err := m.jwtManager.ValidateToken(parts[1])
        if err != nil {
            response.Unauthorized(c, "Token无效或已过期")
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

#### 密码加密
```go
func HashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

func VerifyPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
```

## 前端实现

### 文件结构

```
frontend/
├── src/
│   ├── api/
│   │   └── user-api.js                   # 用户API客户端
│   ├── stores/
│   │   └── user.js                       # 用户状态管理
│   ├── views/
│   │   ├── LoginView.vue                 # 登录/注册页面
│   │   └── ProfileView.vue               # 个人中心页面
│   └── router/
│       └── index.js                      # 路由配置（含认证守卫）
└── .env.development                      # 开发环境变量
```

### 用户状态管理

```javascript
// stores/user.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const currentUser = ref(JSON.parse(localStorage.getItem('currentUser') || 'null'))
  const isAuthenticated = computed(() => !!token.value)

  function setAuth(tokenValue, user) {
    token.value = tokenValue
    currentUser.value = user
    localStorage.setItem('token', tokenValue)
    localStorage.setItem('currentUser', JSON.stringify(user))
  }

  function clearAuth() {
    token.value = ''
    currentUser.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
  }

  return { token, currentUser, isAuthenticated, setAuth, clearAuth }
})
```

### 路由守卫

```javascript
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})
```

### API客户端

```javascript
// api/user-api.js
import axios from 'axios'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080',
  timeout: 10000
})

// 请求拦截器：添加token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：处理401
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('currentUser')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data)
  }
)

export const userApi = {
  register: (data) => apiClient.post('/api/v2/auth/register', data),
  login: (data) => apiClient.post('/api/v2/auth/login', data),
  getProfile: () => apiClient.get('/api/v2/users/profile'),
  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
    return Promise.resolve()
  }
}
```

## 使用说明

### 1. 后端启动

```bash
cd backend

# 首次运行需要执行数据库迁移（可选，GORM会自动创建表）
go run scripts/migrate.go ./data/poems.db

# 启动服务
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 2. 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端将在 `http://localhost:5173` 启动。

### 3. 访问页面

- 登录/注册: `http://localhost:5173/login`
- 个人中心: `http://localhost:5173/profile`（需要登录）

## 安全注意事项

1. **JWT密钥**: 生产环境请修改 `backend/api/router.go` 中的 `jwtManager` 密钥
2. **密码复杂度**: 前端要求密码至少6位，生产环境建议增加更多验证
3. **HTTPS**: 生产环境必须使用HTTPS
4. **CORS**: 确保CORS配置正确，避免跨域问题

## 后续扩展

基于当前实现，可以轻松添加以下功能：

1. **邮箱验证**: 注册时发送验证邮件
2. **密码重置**: 忘记密码功能
3. **第三方登录**: 微信、GitHub等OAuth登录
4. **用户权限**: 角色和权限管理
5. **收藏功能**: 使用 user_favorites 表
6. **浏览历史**: 使用 user_history 表

## 测试

### 使用cURL测试API

```bash
# 注册
curl -X POST http://localhost:8080/api/v2/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456","nickname":"测试"}'

# 登录
curl -X POST http://localhost:8080/api/v2/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 获取个人资料（替换TOKEN）
curl -X GET http://localhost:8080/api/v2/users/profile \
  -H "Authorization: Bearer TOKEN"
```

## 相关文件

- [开发架构设计文档](development-architecture-v2.md)
- [数据库设计](database_design.md)
- [API参考](api-reference.md)
