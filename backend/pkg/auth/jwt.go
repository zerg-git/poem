package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager JWT管理器
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken 生成JWT token
func (m *JWTManager) GenerateToken(userID uint, username string) (string, error) {
	// 生成唯一的JTI
	jti := uuid.New().String()

	now := time.Now()
	expiresAt := now.Add(m.tokenDuration)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT token
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken 刷新token
func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查token是否即将过期（30分钟内）
	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return "", errors.New("token is still valid")
	}

	// 生成新token
	return m.GenerateToken(claims.UserID, claims.Username)
}
