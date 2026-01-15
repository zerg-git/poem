package models

import (
	"time"
)

// User 用户表
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	OpenID       string     `gorm:"size:100;unique" json:"open_id,omitempty"`           // 微信OpenID（预留）
	UnionID      string     `gorm:"size:100" json:"union_id,omitempty"`                 // 微信UnionID（预留）
	Username     string     `gorm:"size:50;not null;unique" json:"username"`            // 用户名
	PasswordHash string     `gorm:"size:255;not null" json:"-"`                         // 密码哈希（不返回给前端）
	Nickname     string     `gorm:"size:100" json:"nickname"`                            // 昵称
	AvatarURL    string     `gorm:"size:500" json:"avatar_url"`                          // 头像URL
	Email        string     `gorm:"size:100;unique" json:"email,omitempty"`             // 邮箱
	Phone        string     `gorm:"size:20;unique" json:"phone,omitempty"`              // 手机号
	Gender       int        `gorm:"default:0;comment:0:未知 1:男 2:女" json:"gender"`  // 性别
	BirthDate    *time.Time `json:"birth_date,omitempty"`                               // 生日
	Province     string     `gorm:"size:50" json:"province,omitempty"`                  // 省份
	City         string     `gorm:"size:50" json:"city,omitempty"`                      // 城市
	Level        int        `gorm:"default:1;comment:用户等级" json:"level"`            // 等级
	Experience   int        `gorm:"default:0;comment:经验值" json:"experience"`          // 经验值
	Coins        int        `gorm:"default:0;comment:金币" json:"coins"`                // 金币
	VIPLevel     int        `gorm:"default:0;comment:VIP等级" json:"vip_level"`        // VIP等级
	VIPExpireAt  *time.Time `json:"vip_expire_at,omitempty"`                            // VIP过期时间
	Status       int        `gorm:"default:1;comment:0:禁用 1:正常" json:"status"`      // 状态
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`                            // 最后登录时间
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserFavorite 用户收藏表
type UserFavorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_target" json:"user_id"`
	TargetID  uint      `gorm:"not null;index:idx_user_target" json:"target_id"`    // 诗词ID或作者ID
	TargetType string    `gorm:"size:20;not null;index:idx_user_target" json:"target_type"` // poem / author
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserFavorite) TableName() string {
	return "user_favorites"
}

// UserHistory 用户浏览历史表
type UserHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	TargetID  uint      `gorm:"not null" json:"target_id"`              // 诗词ID或作者ID
	TargetType string    `gorm:"size:20;not null" json:"target_type"`   // poem / author
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (UserHistory) TableName() string {
	return "user_history"
}
