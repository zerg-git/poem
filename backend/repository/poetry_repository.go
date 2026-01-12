package repository

import (
	"fmt"
	"math/rand/v2"
	"poem/backend/models"
	"strings"
	"sync"
)

// PoetryRepository 诗词数据仓库
type PoetryRepository struct {
	loader   *JSONLoader
	poems    map[string][]models.Poem // category -> poems
	authors  map[string]*models.Author
	mu       sync.RWMutex
}

// NewPoetryRepository 创建诗词仓库
func NewPoetryRepository(dataPath string) *PoetryRepository {
	return &PoetryRepository{
		loader:  NewJSONLoader(dataPath),
		poems:   make(map[string][]models.Poem),
		authors: make(map[string]*models.Author),
	}
}

// GetPoems 获取诗词列表（分页）
func (r *PoetryRepository) GetPoems(page, pageSize int, category string) (models.PoemCollection, error) {
	var allPoems []models.Poem
	var err error

	if category != "" {
		allPoems, err = r.loadCategory(category)
	} else {
		// 加载所有分类
		categories := r.loader.GetCategories()
		for _, cat := range categories {
			poems, loadErr := r.loadCategory(cat.Path)
			if loadErr != nil {
				continue
			}
			// 创建副本以避免共享底层数组
			allPoems = append(allPoems, poems...)
		}
	}

	if err != nil {
		return models.PoemCollection{}, err
	}

	total := len(allPoems)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	// 创建结果切片的副本
	result := make([]models.Poem, end-start)
	copy(result, allPoems[start:end])

	return models.PoemCollection{
		Poems:      result,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	}, nil
}

// GetPoemByID 根据ID获取诗词
func (r *PoetryRepository) GetPoemByID(id string) (*models.Poem, error) {
	categories := r.loader.GetCategories()

	for _, cat := range categories {
		poems, err := r.loadCategory(cat.Path)
		if err != nil {
			continue
		}

		for i := range poems {
			if poems[i].ID == id {
				return &poems[i], nil
			}
		}
	}

	return nil, fmt.Errorf("poem not found")
}

// GetRandomPoems 获取随机诗词
func (r *PoetryRepository) GetRandomPoems(count int, category string) ([]models.Poem, error) {
	var allPoems []models.Poem

	if category != "" {
		poems, err := r.loadCategory(category)
		if err != nil {
			return nil, err
		}
		allPoems = poems
	} else {
		// 从唐诗中获取
		poems, err := r.loadCategory("全唐诗")
		if err == nil && len(poems) > 0 {
			allPoems = poems
		}
	}

	if len(allPoems) == 0 {
		return []models.Poem{}, nil
	}

	// 简单随机选取
	if count > len(allPoems) {
		count = len(allPoems)
	}

	result := make([]models.Poem, 0, count)
	used := make(map[int]bool)

	for len(result) < count {
		idx := rand.IntN(len(allPoems))
		if used[idx] {
			continue
		}
		used[idx] = true
		result = append(result, allPoems[idx])
	}

	return result, nil
}

// GetPoemsByAuthor 根据作者获取诗词
func (r *PoetryRepository) GetPoemsByAuthor(authorName string, page, pageSize int) (models.PoemCollection, error) {
	categories := r.loader.GetCategories()
	var allPoems []models.Poem

	for _, cat := range categories {
		poems, err := r.loadCategory(cat.Path)
		if err != nil {
			continue
		}

		for _, poem := range poems {
			if poem.Author == authorName {
				allPoems = append(allPoems, poem)
			}
		}
	}

	total := len(allPoems)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	// 创建结果切片的副本
	result := make([]models.Poem, end-start)
	copy(result, allPoems[start:end])

	return models.PoemCollection{
		Poems:      result,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	}, nil
}

// GetAuthors 获取作者列表
func (r *PoetryRepository) GetAuthors(page, pageSize int, dynasty string) (models.AuthorCollection, error) {
	// 从诗词中提取作者
	authorMap := make(map[string]*models.Author)

	categories := r.loader.GetCategories()
	for _, cat := range categories {
		if dynasty != "" && cat.DynastyID != dynasty {
			continue
		}

		poems, err := r.loadCategory(cat.Path)
		if err != nil {
			continue
		}

		for _, poem := range poems {
			if poem.Author == "" {
				continue
			}

			if _, exists := authorMap[poem.Author]; !exists {
				authorMap[poem.Author] = &models.Author{
					ID:      poem.Author,
					Name:    poem.Author,
					Dynasty: r.getDynastyName(poem.DynastyID),
				}
			}
			authorMap[poem.Author].PoemCount++
		}
	}

	// 转换为数组
	var authors []models.Author
	for _, author := range authorMap {
		authors = append(authors, *author)
	}

	total := len(authors)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	// 创建结果切片的副本
	result := make([]models.Author, end-start)
	copy(result, authors[start:end])

	return models.AuthorCollection{
		Authors:    result,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	}, nil
}

// GetAuthorByName 根据名称获取作者
func (r *PoetryRepository) GetAuthorByName(name string) (*models.Author, error) {
	authors, err := r.GetAuthors(1, 1000, "")
	if err != nil {
		return nil, err
	}

	for _, author := range authors.Authors {
		if author.Name == name {
			return &author, nil
		}
	}

	return nil, fmt.Errorf("author not found")
}

// Search 搜索诗词
func (r *PoetryRepository) Search(query string, page, pageSize int) (models.SearchResponse, error) {
	categories := r.loader.GetCategories()
	var results []models.Poem

	searchLower := fmt.Sprintf("%s", query)

	for _, cat := range categories {
		poems, err := r.loadCategory(cat.Path)
		if err != nil {
			continue
		}

		for _, poem := range poems {
			// 搜索标题
			if contains(poem.Title, searchLower) {
				results = append(results, poem)
				continue
			}
			// 搜索作者
			if contains(poem.Author, searchLower) {
				results = append(results, poem)
				continue
			}
			// 搜索内容
			for _, paragraph := range poem.Paragraphs {
				if contains(paragraph, searchLower) {
					results = append(results, poem)
					break
				}
			}
		}
	}

	total := len(results)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	// 创建结果切片的副本
	result := make([]models.Poem, end-start)
	copy(result, results[start:end])

	return models.SearchResponse{
		Poems:      result,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
		Query:      query,
	}, nil
}

// loadCategory 加载分类诗词（带缓存）
func (r *PoetryRepository) loadCategory(categoryPath string) ([]models.Poem, error) {
	r.mu.RLock()
	if poems, exists := r.poems[categoryPath]; exists {
		r.mu.RUnlock()
		return poems, nil
	}
	r.mu.RUnlock()

	poems, err := r.loader.LoadCategoryFiles(categoryPath)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.poems[categoryPath] = poems
	r.mu.Unlock()

	return poems, nil
}

// getDynastyName 获取朝代名称
func (r *PoetryRepository) getDynastyName(dynastyID string) string {
	switch dynastyID {
	case "tang":
		return "唐"
	case "song":
		return "宋"
	case "yuan":
		return "元"
	case "wudai":
		return "五代"
	case "ming":
		return "明"
	case "qing":
		return "清"
	default:
		return ""
	}
}

// contains 字符串包含检查
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
