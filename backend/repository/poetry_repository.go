package repository

import (
	"poem/backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PoetryRepository 诗词数据仓库
type PoetryRepository struct {
	db *gorm.DB
}

// NewPoetryRepository 创建诗词仓库
func NewPoetryRepository(dbPath string) (*PoetryRepository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return &PoetryRepository{db: db}, nil
}

// GetPoems 获取诗词列表（分页）
func (r *PoetryRepository) GetPoems(page, pageSize int, categoryName string) (models.PoemCollection, error) {
	var works []models.Work
	var total int64

	query := r.db.Model(&models.Work{}).Preload("Author").Preload("Category")

	if categoryName != "" {
		// Join categories table to filter by category name
		query = query.Joins("JOIN categories ON categories.id = works.category_id").
			Where("categories.name = ? OR categories.display_name = ?", categoryName, categoryName)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("works.id asc").Offset(offset).Limit(pageSize).Find(&works).Error
	if err != nil {
		return models.PoemCollection{}, err
	}

	return models.PoemCollection{
		Works:      works,
		Total:      int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetPoemByID 根据ID获取诗词
func (r *PoetryRepository) GetPoemByID(id string) (*models.Work, error) {
	var work models.Work
	// Try searching by OriginalID first, then Primary Key ID if it's numeric
	err := r.db.Preload("Author").Preload("Category").Preload("Comments").
		Where("original_id = ?", id).First(&work).Error

	if err != nil {
		// If not found by original_id, try by primary key
		err = r.db.Preload("Author").Preload("Category").Preload("Comments").
			First(&work, "id = ?", id).Error
		if err != nil {
			return nil, err
		}
	}
	return &work, nil
}

// GetRandomPoems 获取随机诗词
func (r *PoetryRepository) GetRandomPoems(count int, categoryName string) ([]models.Work, error) {
	var works []models.Work

	query := r.db.Model(&models.Work{}).Preload("Author").Preload("Category")

	if categoryName != "" {
		query = query.Joins("JOIN categories ON categories.id = works.category_id").
			Where("categories.name = ? OR categories.display_name = ?", categoryName, categoryName)
	}

	// SQLite 随机排序
	err := query.Order("RANDOM()").Limit(count).Find(&works).Error
	if err != nil {
		return nil, err
	}

	return works, nil
}

// GetPoemsByAuthor 根据作者获取诗词
func (r *PoetryRepository) GetPoemsByAuthor(authorName string, page, pageSize int) (models.PoemCollection, error) {
	var works []models.Work
	var total int64

	// Find Author first
	var author models.Author
	if err := r.db.Where("name = ?", authorName).First(&author).Error; err != nil {
		return models.PoemCollection{}, err
	}

	query := r.db.Model(&models.Work{}).Preload("Author").Preload("Category").Where("author_id = ?", author.ID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&works).Error
	if err != nil {
		return models.PoemCollection{}, err
	}

	return models.PoemCollection{
		Works:      works,
		Total:      int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetAuthors 获取作者列表
func (r *PoetryRepository) GetAuthors(page, pageSize int, dynasty string) (models.AuthorCollection, error) {
	var authors []models.Author
	var total int64

	query := r.db.Model(&models.Author{})
	if dynasty != "" {
		query = query.Where("dynasty = ?", dynasty)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&authors).Error
	if err != nil {
		return models.AuthorCollection{}, err
	}

	return models.AuthorCollection{
		Authors:    authors,
		Total:      int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetAuthorByName 根据名称获取作者
func (r *PoetryRepository) GetAuthorByName(name string) (*models.Author, error) {
	var author models.Author
	err := r.db.Where("name = ?", name).First(&author).Error
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// Search 搜索诗词
func (r *PoetryRepository) Search(queryStr string, page, pageSize int) (models.SearchResponse, error) {
	var works []models.Work
	var total int64

	// 使用 LIKE 进行模糊搜索
	likeStr := "%" + queryStr + "%"

	// Join with Author to search by author name as well
	query := r.db.Model(&models.Work{}).Preload("Author").Preload("Category").
		Joins("LEFT JOIN authors ON authors.id = works.author_id").
		Where("works.title LIKE ? OR works.content LIKE ? OR authors.name LIKE ?", likeStr, likeStr, likeStr)

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&works).Error
	if err != nil {
		return models.SearchResponse{}, err
	}

	return models.SearchResponse{
		Works:      works,
		Total:      int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		Query:      queryStr,
	}, nil
}

// GetCategories 获取所有分类
func (r *PoetryRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}
