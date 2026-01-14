package services

import (
	"poem/backend/models"
	"poem/backend/repository"
)

// PoetryService 诗词服务
type PoetryService struct {
	repo *repository.PoetryRepository
}

// NewPoetryService 创建诗词服务
func NewPoetryService(repo *repository.PoetryRepository) *PoetryService {
	return &PoetryService{repo: repo}
}

// GetPoems 获取诗词列表
func (s *PoetryService) GetPoems(page, pageSize int, categoryName string) (models.PoemCollection, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetPoems(page, pageSize, categoryName)
}

// GetPoemByID 获取单首诗词
func (s *PoetryService) GetPoemByID(id string) (*models.Work, error) {
	return s.repo.GetPoemByID(id)
}

// GetRandomPoems 获取随机诗词
func (s *PoetryService) GetRandomPoems(count int, categoryName string) ([]models.Work, error) {
	if count < 1 {
		count = 1
	}
	if count > 10 {
		count = 10
	}

	return s.repo.GetRandomPoems(count, categoryName)
}

// GetPoemsByAuthor 获取作者的诗词
func (s *PoetryService) GetPoemsByAuthor(authorName string, page, pageSize int) (models.PoemCollection, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetPoemsByAuthor(authorName, page, pageSize)
}

// GetAuthors 获取作者列表
func (s *PoetryService) GetAuthors(page, pageSize int, dynasty string) (models.AuthorCollection, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetAuthors(page, pageSize, dynasty)
}

// GetAuthorByName 获取作者详情
func (s *PoetryService) GetAuthorByName(name string) (*models.Author, error) {
	return s.repo.GetAuthorByName(name)
}

// Search 搜索
func (s *PoetryService) Search(query string, page, pageSize int) (models.SearchResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.Search(query, page, pageSize)
}

// GetCategories 获取分类列表
func (s *PoetryService) GetCategories() []models.Category {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return []models.Category{}
	}
	return categories
}
