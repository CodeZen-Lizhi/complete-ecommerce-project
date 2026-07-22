package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const demoUserIDHeader = "X-Demo-User-ID"

var safeIdentifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

type demoIdentityContextKey struct{}

// demoIdentityMiddleware 从受控开发 Header 读取身份，并写入标准 request Context。
func demoIdentityMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		userID := strings.TrimSpace(ginContext.GetHeader(demoUserIDHeader))
		if err := validateIdentifier("X-Demo-User-ID", userID); err != nil {
			writeJSONError(ginContext, 401, "unauthenticated", "缺少或无效的开发身份")
			ginContext.Abort()
			return
		}
		ginContext.Request = ginContext.Request.WithContext(withDemoIdentity(ginContext.Request.Context(), userID))
		ginContext.Next()
	}
}

// withDemoIdentity 将已验证的开发身份附加到 Context，避免业务层读取 HTTP Header。
func withDemoIdentity(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, demoIdentityContextKey{}, userID)
}

// demoUserIDFromContext 提取已验证身份；缺失或类型异常均返回稳定错误。
func demoUserIDFromContext(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("请求 Context 不能为空")
	}
	userID, ok := ctx.Value(demoIdentityContextKey{}).(string)
	if !ok || validateIdentifier("X-Demo-User-ID", userID) != nil {
		return "", fmt.Errorf("开发身份不存在或无效")
	}
	return userID, nil
}

// validateIdentifier 校验用户和会话标识，防止 Redis Key 边界被任意字符污染。
func validateIdentifier(name string, value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || len(trimmed) > maximumIdentifierLength || !safeIdentifierPattern.MatchString(trimmed) {
		return fmt.Errorf("%s 必须是 %d 字符以内的安全标识", name, maximumIdentifierLength)
	}
	return nil
}
