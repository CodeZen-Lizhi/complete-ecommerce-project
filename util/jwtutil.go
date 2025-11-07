package util

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

// UserClaims 简化的JWT声明结构，仅包含必要的用户标识
type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

const secret = "GetRandom"
const expires = 24 * time.Hour

// GenerateToken 生成包含用户ID的Token
func GenerateToken(userID string) (string, error) {
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
