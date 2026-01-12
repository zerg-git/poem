package models

// Dynasty 朝代
type Dynasty struct {
	ID          string `json:"id"`
	Name        string `json:"name"`        // 唐代、宋代
	NameEn      string `json:"name_en"`     // Tang, Song
	Description string `json:"description"`
	Period      string `json:"period"`      // 618-907
	SortOrder   int    `json:"sort_order"`
	PoemCount   int    `json:"poem_count,omitempty"`
}

// Category 诗词分类
type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`        // 唐诗、宋词
	NameEn      string `json:"name_en"`     // Tang Poetry
	Description string `json:"description"`
	DynastyID   string `json:"dynasty_id"`
	Path        string `json:"path"`        // 文件路径
}

// Author 作者
type Author struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Dynasty     string `json:"dynasty"`
	Description string `json:"desc"`
	ShortDesc   string `json:"short_description,omitempty"`
	PoemCount   int    `json:"poem_count"`
}

// Poem 诗词
type Poem struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"` // 诗词正文

	// 可选字段
	Content   []string `json:"content,omitempty"`   // 备用字段（诗经）
	Para      []string `json:"para,omitempty"`      // 备用字段（纳兰性德）
	Rhythmic  string   `json:"rhythmic,omitempty"`  // 词牌名（宋词）
	Dynasty   string   `json:"dynasty,omitempty"`   // 朝代（元曲）
	Tags      []string `json:"tags,omitempty"`      // 标签
	Chapter   string   `json:"chapter,omitempty"`   // 章节（诗经）
	Section   string   `json:"section,omitempty"`   // 章节（诗经）

	// 元数据
	CategoryID string `json:"category_id,omitempty"`
	DynastyID  string `json:"dynasty_id,omitempty"`
	AuthorID   string `json:"author_id,omitempty"`
}

// PoemCollection 诗词集合（分页）
type PoemCollection struct {
	Poems      []Poem `json:"poems"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// AuthorCollection 作者集合（分页）
type AuthorCollection struct {
	Authors    []Author `json:"authors"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Poems      []Poem   `json:"poems"`
	Authors    []Author `json:"authors,omitempty"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
	Query      string   `json:"query"`
	DurationMs int64    `json:"duration_ms"`
}

// APIResponse 统一API响应
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}
