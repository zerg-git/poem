package v2

import (
	"poem/backend/api/middleware"
	"poem/backend/pkg/response"
	"poem/backend/services/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *user.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		if err == user.ErrUserAlreadyExists {
			response.Error(c, 409, "用户名已存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		if err == user.ErrInvalidCredentials {
			response.Error(c, 401, "用户名或密码错误")
			return
		}
		if err == user.ErrUserDisabled {
			response.Error(c, 403, "用户已被禁用")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, profile)
}

// RefreshToken 刷新token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	newToken, err := h.userService.RefreshToken(req.Token)
	if err != nil {
		response.Unauthorized(c, "Token刷新失败")
		return
	}

	response.Success(c, gin.H{
		"token": newToken,
	})
}

// GetProfileByID 根据ID获取用户资料（公开接口）
func (h *UserHandler) GetProfileByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 404, "用户不存在")
		return
	}

	// 公开接口只返回部分信息
	publicProfile := gin.H{
		"id":         profile.ID,
		"username":   profile.Username,
		"nickname":   profile.Nickname,
		"avatar_url": profile.AvatarURL,
		"level":      profile.Level,
	}

	response.Success(c, publicProfile)
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req user.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	profile, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if err.Error() == "邮箱已被使用" {
			response.Error(c, 409, "邮箱已被使用")
			return
		}
		if err.Error() == "手机号已被使用" {
			response.Error(c, 409, "手机号已被使用")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, profile)
}
