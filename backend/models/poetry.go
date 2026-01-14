package models

// Dynasty (Legacy support or for static config if needed, though now Author has Dynasty field)
type Dynasty struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NameEn      string `json:"name_en"`
	Description string `json:"description"`
	Period      string `json:"period"`
	SortOrder   int    `json:"sort_order"`
}

// PoemCollection 诗词集合（分页）
type PoemCollection struct {
	Works      []Work `json:"works"`
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
	Works      []Work   `json:"works"`
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
