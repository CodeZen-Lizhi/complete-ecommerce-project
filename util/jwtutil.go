package util

import (
	"ecommerce/internal/config"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtIssuer            = "ecommerce-api"
	jwtAudience          = "ecommerce-client"
	jwtSubject           = "access-token"
	defaultWeakJWTSecret = "change_me_to_a_random_secret"
	minJWTSecretLength   = 32
)

// UserClaims 简化的JWT声明结构，仅包含必要的用户标识
type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// jwtConfig 读取并校验JWT配置。
func jwtConfig() (string, time.Duration, error) {
	if config.Cfg == nil {
		return "", 0, errors.New("配置未初始化")
	}
	secret := strings.TrimSpace(config.Cfg.Jwt.Secret)
	if secret == "" {
		return "", 0, errors.New("jwt.secret 不能为空")
	}
	if secret == defaultWeakJWTSecret {
		return "", 0, errors.New("jwt.secret 不能使用默认弱密钥，请更换为高强度随机字符串")
	}
	if len(secret) < minJWTSecretLength {
		return "", 0, fmt.Errorf("jwt.secret 长度不能小于 %d", minJWTSecretLength)
	}
	if config.Cfg.Jwt.ExpiresHours <= 0 {
		return "", 0, errors.New("jwt.expires_hours 必须大于0")
	}
	expires := time.Duration(config.Cfg.Jwt.ExpiresHours) * time.Hour
	return secret, expires, nil
}

// GenerateToken 生成包含用户ID的Token
func GenerateToken(userID int64) (string, error) {
	secret, expires, err := jwtConfig()
	if err != nil {
		return "", err
	}
	// 设置过期时间
	now := time.Now()
	expireAt := now.Add(expires)

	// 构建声明
	claims := UserClaims{
		UserID: strconv.FormatInt(userID, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   jwtSubject,
			Audience:  jwt.ClaimStrings{jwtAudience},
			ExpiresAt: jwt.NewNumericDate(expireAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// 生成并签名Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 从Token中获取用户ID
func ParseToken(tokenString string) (int64, error) {
	secret, _, err := jwtConfig()
	if err != nil {
		return 0, err
	}
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return 0, errors.New("token不能为空")
	}
	// 解析Token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuedAt(),
	)

	if err != nil {
		return 0, err
	}

	// 提取用户ID
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// 兼容旧 token：旧 token 可能没有以下标准字段
		if claims.Issuer != "" && claims.Issuer != jwtIssuer {
			return 0, errors.New("token签发者无效")
		}
		if claims.Subject != "" && claims.Subject != jwtSubject {
			return 0, errors.New("token主题无效")
		}
		if len(claims.Audience) > 0 {
			validAudience := false
			for _, audience := range claims.Audience {
				if audience == jwtAudience {
					validAudience = true
					break
				}
			}
			if !validAudience {
				return 0, errors.New("token受众无效")
			}
		}
		// 将 string 类型的 UserID 转换为 int64
		userID, err := strconv.ParseInt(claims.UserID, 10, 64)
		if err != nil {
			return 0, errors.New("用户ID格式错误")
		}
		return userID, nil
	}

	return 0, errors.New("无效的token")
}
