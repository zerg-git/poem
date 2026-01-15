-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    open_id VARCHAR(100) UNIQUE,
    union_id VARCHAR(100),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar_url VARCHAR(500),
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    gender INTEGER DEFAULT 0 CHECK(gender IN (0, 1, 2)),
    birth_date DATE,
    province VARCHAR(50),
    city VARCHAR(50),
    level INTEGER DEFAULT 1,
    experience INTEGER DEFAULT 0,
    coins INTEGER DEFAULT 0,
    vip_level INTEGER DEFAULT 0,
    vip_expire_at DATETIME,
    status INTEGER DEFAULT 1 CHECK(status IN (0, 1)),
    last_login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户收藏表
CREATE TABLE IF NOT EXISTS user_favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    target_type VARCHAR(20) NOT NULL CHECK(target_type IN ('poem', 'author')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, target_id, target_type)
);

-- 创建收藏表索引
CREATE INDEX IF NOT EXISTS idx_favorites_user ON user_favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_target ON user_favorites(target_id, target_type);

-- 用户浏览历史表
CREATE TABLE IF NOT EXISTS user_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    target_type VARCHAR(20) NOT NULL CHECK(target_type IN ('poem', 'author')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建历史表索引
CREATE INDEX IF NOT EXISTS idx_history_user ON user_history(user_id);
CREATE INDEX IF NOT EXISTS idx_history_created ON user_history(created_at DESC);
