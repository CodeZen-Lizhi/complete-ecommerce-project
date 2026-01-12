package util

import (
	"ecommerce/internal/config"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims 简化的JWT声明结构，仅包含必要的用户标识
type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func jwtConfig() (string, time.Duration, error) {
	if config.Cfg == nil {
		return "", 0, errors.New("配置未初始化")
	}
	secret := config.Cfg.Jwt.Secret
	if secret == "" {
		return "", 0, errors.New("jwt.secret 不能为空")
	}
	if config.Cfg.Jwt.ExpiresHours <= 0 {
		return "", 0, errors.New("jwt.expires_hours 必须大于0")
	}
	expires := time.Duration(config.Cfg.Jwt.ExpiresHours) * time.Hour
	return secret, expires, nil
}

// GenerateToken 生成包含用户ID的Token
func GenerateToken(userID string) (string, error) {
	secret, expires, err := jwtConfig()
	if err != nil {
		return "", err
	}
	// 设置过期时间
	expireAt := time.Now().Add(expires)

	// 构建声明
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 生成并签名Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 从Token中获取用户ID
func ParseToken(tokenString string) (uint64, error) {
	secret, _, err := jwtConfig()
	if err != nil {
		return 0, err
	}
	// 解析Token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return 0, err
	}

	// 提取用户ID
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// 将 string 类型的 UserID 转换为 uint64
		userID, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			return 0, errors.New("用户ID格式错误")
		}
		return userID, nil
	}

	return 0, errors.New("无效的token")
}
