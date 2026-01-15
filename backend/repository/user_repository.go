package repository

import (
	"context"
	"errors"
	"poem/backend/models"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *models.User) error
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uint) (*models.User, error)
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	// GetByPhone 根据手机号获取用户
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	// Update 更新用户
	Update(ctx context.Context, user *models.User) error
	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, userID uint) error
	// AddFavorite 添加收藏
	AddFavorite(ctx context.Context, userID, targetID uint, targetType string) error
	// RemoveFavorite 取消收藏
	RemoveFavorite(ctx context.Context, userID, targetID uint, targetType string) error
	// GetFavorites 获取收藏列表
	GetFavorites(ctx context.Context, userID uint, targetType string, page, pageSize int) ([]models.UserFavorite, int64, error)
	// AddHistory 添加浏览历史
	AddHistory(ctx context.Context, userID, targetID uint, targetType string) error
	// GetHistory 获取浏览历史
	GetHistory(ctx context.Context, userID uint, targetType string, page, pageSize int) ([]models.UserHistory, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户Repository
func NewUserRepository(db *gorm.DB) (UserRepository, error) {
	// 自动迁移表结构
	if err := db.AutoMigrate(&models.User{}, &models.UserFavorite{}, &models.UserHistory{}); err != nil {
		return nil, err
	}
	return &userRepository{db: db}, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ? AND status = 1", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", now).Error
}

func (r *userRepository) AddFavorite(ctx context.Context, userID, targetID uint, targetType string) error {
	favorite := &models.UserFavorite{
		UserID:    userID,
		TargetID:  targetID,
		TargetType: targetType,
	}
	// 使用 FirstOrCreate 避免重复
	return r.db.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		FirstOrCreate(favorite).Error
}

func (r *userRepository) RemoveFavorite(ctx context.Context, userID, targetID uint, targetType string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Delete(&models.UserFavorite{}).Error
}

func (r *userRepository) GetFavorites(ctx context.Context, userID uint, targetType string, page, pageSize int) ([]models.UserFavorite, int64, error) {
	var favorites []models.UserFavorite
	var total int64

	query := r.db.WithContext(ctx).Model(&models.UserFavorite{}).Where("user_id = ?", userID)
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&favorites).Error

	return favorites, total, err
}

func (r *userRepository) AddHistory(ctx context.Context, userID, targetID uint, targetType string) error {
	history := &models.UserHistory{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
	}
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *userRepository) GetHistory(ctx context.Context, userID uint, targetType string, page, pageSize int) ([]models.UserHistory, int64, error) {
	var history []models.UserHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.UserHistory{}).Where("user_id = ?", userID)
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&history).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.UserHistory{}, 0, nil
	}

	return history, total, err
}
