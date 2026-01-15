package user

import (
	"context"
	"errors"
	"poem/backend/pkg/auth"
	"poem/backend/repository"
	"poem/backend/models"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("用户已存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserDisabled      = errors.New("用户已被禁用")
)

// UserService 用户服务
type UserService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string     `json:"token"`
	ExpiresAt int64      `json:"expires_at"`
	User      *UserInfo  `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Gender     int    `json:"gender"`
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
	Coins      int    `json:"coins"`
	VIPLevel   int    `json:"vip_level"`
	Status     int    `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// toUserInfo 将User模型转换为UserInfo
func toUserInfo(user *models.User) *UserInfo {
	return &UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Nickname:   user.Nickname,
		AvatarURL:  user.AvatarURL,
		Email:      user.Email,
		Phone:      user.Phone,
		Gender:     user.Gender,
		Level:      user.Level,
		Experience: user.Experience,
		Coins:      user.Coins,
		VIPLevel:   user.VIPLevel,
		Status:     user.Status,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	// 检查邮箱是否已被使用
	if req.Email != "" {
		_, err = s.userRepo.GetByEmail(ctx, req.Email)
		if err == nil {
			return nil, errors.New("邮箱已被使用")
		}
	}

	// 密码哈希
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       1, // 默认正常
	}

	if req.Nickname == "" {
		user.Nickname = req.Username
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 生成token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		User:      toUserInfo(user),
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 验证密码
	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 生成token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		User:      toUserInfo(user),
	}, nil
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(ctx context.Context, userID uint) (*UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return toUserInfo(user), nil
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(tokenString string) (string, error) {
	return s.jwtManager.RefreshToken(tokenString)
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Nickname  string `json:"nickname" binding:"omitempty,max=100"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,max=500"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty,len=11"`
	Gender    int    `json:"gender" binding:"omitempty,min=0,max=2"`
	Province  string `json:"province" binding:"omitempty,max=50"`
	City      string `json:"city" binding:"omitempty,max=50"`
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*UserInfo, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 如果要更新邮箱，检查邮箱是否已被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}

	// 如果要更新手机，检查手机是否已被其他用户使用
	if req.Phone != "" && req.Phone != user.Phone {
		existingUser, err := s.userRepo.GetByPhone(ctx, req.Phone)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("手机号已被使用")
		}
		user.Phone = req.Phone
	}

	// 更新其他字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Gender >= 0 && req.Gender <= 2 {
		user.Gender = req.Gender
	}
	if req.Province != "" {
		user.Province = req.Province
	}
	if req.City != "" {
		user.City = req.City
	}

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return toUserInfo(user), nil
}
