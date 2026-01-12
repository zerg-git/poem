package repository

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"poem/backend/models"
	"sort"
	"strings"
	"sync"
)

// JSONLoader JSON文件加载器
type JSONLoader struct {
	basePath string
	cache    map[string][]models.Poem
	cacheMu  sync.RWMutex
}

// NewJSONLoader 创建JSON加载器
func NewJSONLoader(basePath string) *JSONLoader {
	return &JSONLoader{
		basePath: basePath,
		cache:    make(map[string][]models.Poem),
	}
}

// LoadCategoryFiles 加载指定分类的所有JSON文件
func (jl *JSONLoader) LoadCategoryFiles(categoryPath string) ([]models.Poem, error) {
	fullPath := filepath.Join(jl.basePath, categoryPath)

	// 检查缓存
	jl.cacheMu.RLock()
	if cached, exists := jl.cache[categoryPath]; exists {
		jl.cacheMu.RUnlock()
		return cached, nil
	}
	jl.cacheMu.RUnlock()

	var allPoems []models.Poem

	err := filepath.Walk(fullPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
			return nil
		}

		baseName := filepath.Base(path)

		// 对于全唐诗，只加载 poet.tang. 或 poet.song. 开头的文件
		if categoryPath == "全唐诗" {
			if !strings.HasPrefix(baseName, "poet.tang.") && !strings.HasPrefix(baseName, "poet.song.") {
				return nil
			}
		}

		// 对于宋词，只加载 ci.song. 开头的文件
		if categoryPath == "宋词" {
			if !strings.HasPrefix(baseName, "ci.song.") {
				return nil
			}
		}

		// 对于五代诗词，跳过子目录中的非诗词文件
		if categoryPath == "五代诗词" {
			// 只加载 JSON 文件，跳过 README 等
			if !strings.HasSuffix(baseName, ".json") {
				return nil
			}
		}

		poems, err := jl.loadPoemsFromFile(path, categoryPath)
		if err != nil {
			// 跳过无法解析的文件
			return nil
		}

		allPoems = append(allPoems, poems...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	// 缓存结果
	jl.cacheMu.Lock()
	jl.cache[categoryPath] = allPoems
	jl.cacheMu.Unlock()

	return allPoems, nil
}

// loadPoemsFromFile 从文件加载诗词
func (jl *JSONLoader) loadPoemsFromFile(filePath, categoryPath string) ([]models.Poem, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 尝试解析为诗词数组
	var poems []models.Poem
	if err := json.Unmarshal(data, &poems); err != nil {
		// 尝试解析原始格式
		return jl.loadRawFormat(data, filePath, categoryPath)
	}

	// 添加元数据并过滤无效数据
	categoryID := jl.extractCategoryID(categoryPath)
	dynastyID := jl.extractDynastyID(categoryPath)

	validPoems := make([]models.Poem, 0, len(poems))
	for i := range poems {
		// 处理content字段（诗经等使用content而非paragraphs）
		if len(poems[i].Paragraphs) == 0 && len(poems[i].Content) > 0 {
			poems[i].Paragraphs = poems[i].Content
		}
		// 处理para字段（纳兰性德等使用para）
		if len(poems[i].Paragraphs) == 0 && len(poems[i].Para) > 0 {
			poems[i].Paragraphs = poems[i].Para
		}

		// 对于宋词，如果没有title但有rhythmic（词牌名），则用rhythmic作为标题
		if poems[i].Title == "" && poems[i].Rhythmic != "" {
			poems[i].Title = poems[i].Rhythmic
		}

		// 过滤无效数据：必须有标题和段落内容
		if poems[i].Title == "" || len(poems[i].Paragraphs) == 0 {
			continue
		}

		if poems[i].ID == "" {
			poems[i].ID = jl.generateID(filePath, i)
		}
		if poems[i].CategoryID == "" {
			poems[i].CategoryID = categoryID
		}
		if poems[i].DynastyID == "" {
			poems[i].DynastyID = dynastyID
		}

		validPoems = append(validPoems, poems[i])
	}

	return validPoems, nil
}

// loadRawFormat 加载原始格式的JSON（诗经等特殊格式）
func (jl *JSONLoader) loadRawFormat(data []byte, filePath, categoryPath string) ([]models.Poem, error) {
	var rawPoems []map[string]interface{}
	if err := json.Unmarshal(data, &rawPoems); err != nil {
		return nil, err
	}

	var poems []models.Poem
	categoryID := jl.extractCategoryID(categoryPath)
	dynastyID := jl.extractDynastyID(categoryPath)

	for i, raw := range rawPoems {
		poem := models.Poem{
			ID:         jl.generateID(filePath, i),
			CategoryID: categoryID,
			DynastyID:  dynastyID,
		}

		if title, ok := raw["title"].(string); ok {
			poem.Title = title
		}
		if author, ok := raw["author"].(string); ok {
			poem.Author = author
		}
		if dynasty, ok := raw["dynasty"].(string); ok {
			poem.Dynasty = dynasty
		}
		if rhythmic, ok := raw["rhythmic"].(string); ok {
			poem.Rhythmic = rhythmic
		}
		if chapter, ok := raw["chapter"].(string); ok {
			poem.Chapter = chapter
		}
		if section, ok := raw["section"].(string); ok {
			poem.Section = section
		}

		// 处理paragraphs
		if paragraphs, ok := raw["paragraphs"].([]interface{}); ok {
			for _, p := range paragraphs {
				if str, ok := p.(string); ok {
					poem.Paragraphs = append(poem.Paragraphs, str)
				}
			}
		}

		// 处理content（诗经）
		if len(poem.Paragraphs) == 0 {
			if content, ok := raw["content"].([]interface{}); ok {
				for _, c := range content {
					if str, ok := c.(string); ok {
						poem.Paragraphs = append(poem.Paragraphs, str)
					}
				}
			}
		}

		// 处理para（纳兰性德）
		if len(poem.Paragraphs) == 0 {
			if para, ok := raw["para"].([]interface{}); ok {
				for _, p := range para {
					if str, ok := p.(string); ok {
						poem.Paragraphs = append(poem.Paragraphs, str)
					}
				}
			}
		}

		// 对于宋词，如果没有title但有rhythmic（词牌名），则用rhythmic作为标题
		if poem.Title == "" && poem.Rhythmic != "" {
			poem.Title = poem.Rhythmic
		}

		// 过滤无效数据：必须有标题和段落内容
		if poem.Title == "" || len(poem.Paragraphs) == 0 {
			continue
		}

		poems = append(poems, poem)
	}

	return poems, nil
}

// extractCategoryID 从路径提取分类ID
func (jl *JSONLoader) extractCategoryID(categoryPath string) string {
	parts := strings.Split(categoryPath, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		return strings.ToLower(strings.ReplaceAll(lastPart, "-", ""))
	}
	return categoryPath
}

// extractDynastyID 从路径提取朝代ID
func (jl *JSONLoader) extractDynastyID(categoryPath string) string {
	switch {
	case strings.Contains(categoryPath, "全唐诗") || strings.Contains(categoryPath, "水墨唐诗"):
		return "tang"
	case strings.Contains(categoryPath, "宋词"):
		return "song"
	case strings.Contains(categoryPath, "元曲"):
		return "yuan"
	case strings.Contains(categoryPath, "五代"):
		return "wudai"
	case strings.Contains(categoryPath, "诗经") || strings.Contains(categoryPath, "楚辞"):
		return "preqin"
	case strings.Contains(categoryPath, "纳兰性德"):
		return "qing"
	case strings.Contains(categoryPath, "曹操"):
		return "other"
	case strings.Contains(categoryPath, "幽梦影"):
		return "other"
	default:
		return "other"
	}
}

// generateID 生成唯一ID
func (jl *JSONLoader) generateID(filePath string, index int) string {
	return fmt.Sprintf("%s-%d", filepath.Base(filePath), index)
}

// GetDynasties 获取所有朝代
func (jl *JSONLoader) GetDynasties() []models.Dynasty {
	dynasties := []models.Dynasty{
		{
			ID:          "preqin",
			Name:        "先秦",
			NameEn:      "Pre-Qin",
			Description: "先秦时期，包含诗经、楚辞等经典",
			Period:      "-207 BC",
			SortOrder:   0,
		},
		{
			ID:          "tang",
			Name:        "唐代",
			NameEn:      "Tang Dynasty",
			Description: "唐朝（618-907年），中国历史上最辉煌的朝代之一，诗歌创作达到顶峰",
			Period:      "618-907",
			SortOrder:   1,
		},
		{
			ID:          "wudai",
			Name:        "五代",
			NameEn:      "Five Dynasties",
			Description: "五代十国时期（907-960年）",
			Period:      "907-960",
			SortOrder:   2,
		},
		{
			ID:          "song",
			Name:        "宋代",
			NameEn:      "Song Dynasty",
			Description: "宋朝（960-1279年），词的创作达到顶峰",
			Period:      "960-1279",
			SortOrder:   3,
		},
		{
			ID:          "yuan",
			Name:        "元代",
			NameEn:      "Yuan Dynasty",
			Description: "元朝（1271-1368年），散曲成为主要文学形式",
			Period:      "1271-1368",
			SortOrder:   4,
		},
		{
			ID:          "ming",
			Name:        "明代",
			NameEn:      "Ming Dynasty",
			Description: "明朝（1368-1644年）",
			Period:      "1368-1644",
			SortOrder:   5,
		},
		{
			ID:          "qing",
			Name:        "清代",
			NameEn:      "Qing Dynasty",
			Description: "清朝（1644-1912年），纳兰性德等词人活跃时期",
			Period:      "1644-1912",
			SortOrder:   6,
		},
		{
			ID:          "other",
			Name:        "其他",
			NameEn:      "Other",
			Description: "其他时期或跨朝代作品",
			Period:      "-",
			SortOrder:   99,
		},
	}

	sort.Slice(dynasties, func(i, j int) bool {
		return dynasties[i].SortOrder < dynasties[j].SortOrder
	})

	return dynasties
}

// GetCategories 获取所有分类
func (jl *JSONLoader) GetCategories() []models.Category {
	categories := []models.Category{
		// 先秦
		{
			ID:          "shi-jing",
			Name:        "诗经",
			NameEn:      "Shi Jing",
			Description: "中国古代诗歌开端，收录305篇诗歌",
			DynastyID:   "preqin",
			Path:        "诗经",
		},
		{
			ID:          "chu-ci",
			Name:        "楚辞",
			NameEn:      "Chu Ci",
			Description: "屈原创作的诗歌总集",
			DynastyID:   "preqin",
			Path:        "楚辞",
		},
		// 唐代
		{
			ID:          "tang-shi",
			Name:        "唐诗",
			NameEn:      "Tang Poetry",
			Description: "全唐诗收录唐诗四万八千九百余首",
			DynastyID:   "tang",
			Path:        "全唐诗",
		},
		// 五代
		{
			ID:          "wudai",
			Name:        "五代诗词",
			NameEn:      "Five Dynasties",
			Description: "五代十国时期的诗词作品",
			DynastyID:   "wudai",
			Path:        "五代诗词",
		},
		// 宋代
		{
			ID:          "song-ci",
			Name:        "宋词",
			NameEn:      "Song Lyrics",
			Description: "全宋词收录宋词二万余首",
			DynastyID:   "song",
			Path:        "宋词",
		},
		// 元代
		{
			ID:          "yuan-qu",
			Name:        "元曲",
			NameEn:      "Yuan Qu",
			Description: "元代文学形式，包括散曲和杂剧",
			DynastyID:   "yuan",
			Path:        "元曲",
		},
		// 清代
		{
			ID:          "nalan-xingde",
			Name:        "纳兰性德",
			NameEn:      "Nalan Xingde",
			Description: "清代词人纳兰性德诗集",
			DynastyID:   "qing",
			Path:        "纳兰性德",
		},
		// 其他
		{
			ID:          "shuimo-tangshi",
			Name:        "水墨唐诗",
			NameEn:      "Ink Tang Poetry",
			Description: "水墨风格唐诗精选",
			DynastyID:   "tang",
			Path:        "水墨唐诗",
		},
		{
			ID:          "caocao",
			Name:        "曹操诗集",
			NameEn:      "Cao Cao Poetry",
			Description: "曹操诗歌全集",
			DynastyID:   "other",
			Path:        "曹操诗集",
		},
		{
			ID:          "youmengying",
			Name:        "幽梦影",
			NameEn:      "You Meng Ying",
			Description: "幽梦影",
			DynastyID:   "other",
			Path:        "幽梦影",
		},
	}

	return categories
}
