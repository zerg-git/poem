package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusOK,
		Success:   true,
		Data:      data,
		Timestamp: getCurrentTimestamp(),
	})
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:      httpStatus,
		Success:   false,
		Error:     message,
		Timestamp: getCurrentTimestamp(),
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusOK,
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: getCurrentTimestamp(),
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// BadRequest 错误请求响应
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// InternalError 内部错误响应
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

func getCurrentTimestamp() int64 {
	return 0 // 可以添加时间戳逻辑
}
