package v2

import (
	"poem/backend/api/handlers/v2"
	"poem/backend/api/middleware"
	"poem/backend/pkg/auth"

	"github.com/gin-gonic/gin"
)

// Router v2路由
type Router struct {
	userHandler   *v2.UserHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewRouter 创建v2路由
func NewRouter(
	userHandler *v2.UserHandler,
	jwtManager *auth.JWTManager,
) *Router {
	return &Router{
		userHandler:   userHandler,
		authMiddleware: middleware.NewAuthMiddleware(jwtManager),
	}
}

// SetupRoutes 设置v2路由
func (r *Router) SetupRoutes(rg *gin.RouterGroup) {
	// 公开路由
	public := rg.Group("")
	{
		// 用户认证
		public.POST("/auth/register", r.userHandler.Register)
		public.POST("/auth/login", r.userHandler.Login)
		public.POST("/auth/refresh", r.userHandler.RefreshToken)

		// 公开用户信息
		public.GET("/users/:id", r.userHandler.GetProfileByID)
	}

	// 需要认证的路由
	protected := rg.Group("")
	protected.Use(r.authMiddleware.RequireAuth())
	{
		// 用户信息
		protected.GET("/users/profile", r.userHandler.GetProfile)
		protected.PUT("/users/profile", r.userHandler.UpdateProfile)
	}
}
