package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Category 分类表
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"` // e.g., 'quantangshi'
	DisplayName string    `gorm:"size:100" json:"display_name"`        // e.g., '全唐诗'
	Description string    `gorm:"type:text" json:"description"`
	Works       []Work    `gorm:"foreignKey:CategoryID" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

// Author 作者表
type Author struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null;index:idx_author_dynasty,unique" json:"name"`
	Dynasty   string    `gorm:"size:50;index:idx_author_dynasty,unique" json:"dynasty"`
	Biography string    `gorm:"type:text" json:"biography"`
	Works     []Work    `gorm:"foreignKey:AuthorID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// Work 作品表
type Work struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"index" json:"category_id"`
	AuthorID   uint      `gorm:"index" json:"author_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Author     Author    `gorm:"foreignKey:AuthorID" json:"author"`
	Title      string    `gorm:"size:255;index" json:"title"`
	Rhythmic   string    `gorm:"size:255;index" json:"rhythmic"` // 词牌名/曲牌名
	Volume     string    `gorm:"size:100" json:"volume"`         // 卷
	Section    string    `gorm:"size:100" json:"section"`        // 篇/章
	Content    JSONArr   `gorm:"type:text;not null" json:"content"`
	Prologue   string    `gorm:"type:text" json:"prologue"`
	OriginalID string    `gorm:"size:100" json:"original_id"`
	Comments   []Comment `gorm:"foreignKey:WorkID" json:"comments"`
	CreatedAt  time.Time `json:"created_at"`
}

// Comment 注释/评析表
type Comment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	WorkID         uint      `gorm:"index" json:"work_id"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	Type           string    `gorm:"size:50;default:'note'" json:"type"` // note, comment, translation
	Commenter      string    `gorm:"size:100" json:"commenter"`
	ParagraphIndex int       `json:"paragraph_index"` // Optional
	CreatedAt      time.Time `json:"created_at"`
}

// JSONArr 用于在 SQLite 中存储 JSON 数组
type JSONArr []string

// Value 实现 driver.Valuer 接口
func (j JSONArr) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONArr) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		// 尝试作为字符串处理
		str, ok := value.(string)
		if !ok {
			return errors.New("failed to scan JSONArr")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, j)
}

// TableName overrides
func (Category) TableName() string { return "categories" }
func (Author) TableName() string   { return "authors" }
func (Work) TableName() string     { return "works" }
func (Comment) TableName() string  { return "comments" }
