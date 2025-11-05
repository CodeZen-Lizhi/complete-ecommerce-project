package middleware

import (
	"context"
	"ecommerce/internal/logger"
	"ecommerce/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

// CorsMiddleware 处理跨域请求
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Role")
		// 允许暴露的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		// 允许携带cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// AuthMiddleware 验证用户是否登录
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		// 1. 从请求头获取Authorization
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"msg": "缺少Authorization头"})
			c.Abort()
			return
		}
		// 2. 剥离Bearer前缀（格式：Bearer <token>）
		var tokenString string
		_, err := fmt.Sscanf(token, "Bearer %s", &tokenString)
		if err != nil || tokenString == "" {
			c.JSON(401, gin.H{"msg": "Authorization格式错误（应为Bearer <token>）"})
			c.Abort()
			return
		}
		// 3. 调用工具类解析令牌
		userID, err := util.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权访问，请先登录",
			})
			return
		}
		// 可以在这里解析token获取用户信息并设置到上下文
		c.Set(util.CurrentUserId, userID)
		c.Next()
	}
}

// AdminAuthMiddleware 验证是否为管理员
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Role头
		role := c.GetHeader("Role")

		// 验证是否为管理员
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "没有管理员权限，无法访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
	//todo 写一个统计方法用时的中间件
}

// RequestID 在 Gin 中间件中生成 request_id 链路追踪
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用客户端传入的 X-Request-ID
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set("X-Request-ID", rid)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "request_id", rid))

		// 日志中加入 request_id
		c.Next()
	}
}

// GinLogger 优化后的请求日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 初始化日志字段（预分配空间减少内存分配）
		fields := make([]any, 0, 8)

		// 2. 记录请求开始时间和基本信息
		startTime := time.Now()
		req := c.Request
		fields = append(fields,
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.String("request_id", c.GetString("X-Request-ID")),
		)

		// 3. 使用自定义writer捕获响应大小
		w := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = w

		// 4. 处理请求
		c.Next()

		// 5. 补充响应信息
		latency := time.Since(startTime)
		fields = append(fields,
			slog.Int("status_code", c.Writer.Status()),
			slog.Duration("latency", latency),
			slog.Int("response_size", w.size),
		)

		// 6. 记录错误信息（如果有）
		if len(c.Errors) > 0 {
			fields = append(fields, slog.String("error", c.Errors.String()))
		}
		// 7. 根据状态码选择日志级别
		log := logger.GetLogger()
		switch {
		case c.Writer.Status() >= http.StatusInternalServerError:
			log.ErrorContext(c.Request.Context(), "请求处理异常", fields...)
		case c.Writer.Status() >= http.StatusBadRequest:
			log.WarnContext(c.Request.Context(), "客户端请求错误", fields...)
		case latency > 500*time.Millisecond: // 慢请求告警
			fields = append(fields, slog.Bool("slow_request", true))
			log.WarnContext(c.Request.Context(), "请求处理耗时过长", fields...)
		default:
			log.InfoContext(c.Request.Context(), "请求处理完成", fields...)
		}
	}
}

// 辅助类型：用于捕获响应大小（保留原有功能）
type responseWriter struct {
	gin.ResponseWriter
	size int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *responseWriter) WriteString(s string) (int, error) {
	n, err := w.ResponseWriter.WriteString(s)
	w.size += n
	return n, err
}

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var err error
			if panicVal := recover(); panicVal != nil {
				switch val := panicVal.(type) {
				case error:
					err = fmt.Errorf("%w", val)
				default:
					err = fmt.Errorf("%v", val)
				}
				//输出结构化日志（无!BADKEY，符合slog规范）
				logPanic(c, logger.GetLogger(), err)
				//返回标准化500响应
				sendErrorResponse(c, c.GetString("X-Request-ID"))
			}
		}()
		c.Next()
	}
}

// logPanic 输出Panic日志
func logPanic(c *gin.Context, logger *slog.Logger, err error) {
	logger.ErrorContext(
		context.Background(),
		"stack", fmt.Sprintf("错误: %v", err),
		"request_id", c.GetString("X-Request-ID"),
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
	)
}

// sendErrorResponse 返回标准化JSON响应（前后端统一格式）
func sendErrorResponse(c *gin.Context, requestID string) {
	// 避免重复响应（若业务Handler已写响应，跳过）
	if c.Writer.Written() {
		return
	}

	c.AbortWithStatusJSON(
		http.StatusInternalServerError,
		gin.H{
			"code":       http.StatusInternalServerError,
			"msg":        "服务器内部错误，请联系技术支持",
			"request_id": requestID,
		},
	)
}
