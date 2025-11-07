package middleware

import (
	"bufio"
	"context"
	"ecommerce/container"
	"ecommerce/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
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

// RequestIdMiddleware 在 Gin 中间件中生成 request_id 链路追踪
func RequestIdMiddleware() gin.HandlerFunc {
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

// InjectContainerMiddleware 将容器注入 Gin Context
func InjectContainerMiddleware(ctn *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(container.ContainerKey, ctn)
		c.Next()
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

// GinRecovery 是一个能精确定位错误的 Panic 恢复中间件
func GinRecovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 1. 将 panic 转换为 error
				var err error
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", v)
				}

				// 2. 捕获并过滤堆栈信息
				stack := debug.Stack()
				filteredStack := filterStack(string(stack))

				// 3. 记录包含精简堆栈的结构化日志
				requestID := c.GetString("X-Request-ID")
				logger.ErrorContext(
					c.Request.Context(),
					"Panic Recovered",
					slog.Any("error", err),
					slog.String("request_id", requestID),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("stacktrace", filteredStack), // 关键：记录过滤后的堆栈
				)

				// 4. 返回标准化的 500 错误响应
				sendErrorResponse(c, requestID)
			}
		}()
		c.Next()
	}
}

// filterStack 通过分析堆栈，精准定位到触发 panic 的业务代码行
func filterStack(stack string) string {
	projectModulePath := findModulePath()
	if projectModulePath == "" {
		return stack // 如果找不到模块路径，返回原始堆栈
	}

	var relevantStack []string
	lines := strings.Split(stack, "\n")

	// 保留 goroutine 状态行
	if len(lines) > 0 {
		relevantStack = append(relevantStack, lines[0])
	}

	// 从堆栈中找到 panic 的直接原因
	// 策略：在原始堆栈中，panic 的源头通常紧跟在 go runtime 的 panic() 函数帧之后
	var panicFrameFound bool
	for i := 1; i < len(lines)-1; i += 2 {
		functionLine := lines[i]
		fileLine := lines[i+1]

		// 标记 panic() 函数帧的位置
		if strings.Contains(functionLine, "panic(") && strings.Contains(fileLine, "runtime/panic.go") {
			panicFrameFound = true
			continue
		}

		// 在找到 panic 帧后，第一个不属于中间件目录且属于我们自己项目的帧就是源头
		isMiddleware := strings.Contains(fileLine, "/middleware/") // 关键：过滤掉所有中间件目录下的文件
		isProjectCode := strings.Contains(fileLine, projectModulePath)

		if panicFrameFound && isProjectCode && !isMiddleware {
			relevantStack = append(relevantStack, functionLine)
			relevantStack = append(relevantStack, fileLine)
			// 成功找到，跳出循环
			break
		}
	}

	// 如果新策略未找到（例如 panic 直接发生在中间件中），则返回原始堆栈以防信息丢失
	if len(relevantStack) <= 1 {
		return stack
	}

	return strings.Join(relevantStack, "\n")
}

// findModulePath 从 go.mod 文件中读取模块路径
func findModulePath() string {
	file, err := os.Open("go.mod")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

// GinLogger 是一个高度优化的结构化日志中间件
func GinLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		req := c.Request
		requestID := c.GetString("X-Request-ID") // 假设由其他中间件设置

		// 使用自定义 writer 捕获响应大小
		w := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// 使用 slog.Attr 组织日志字段，更具可读性
		attrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.Int("status_code", statusCode),
			slog.Duration("latency", latency),
			slog.Int("response_size", w.size),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		// 根据状态码和延迟选择日志级别
		msg := "请求处理完成"
		level := slog.LevelInfo
		switch {
		case statusCode >= http.StatusInternalServerError:
			msg = "请求处理异常"
			level = slog.LevelError
		case statusCode >= http.StatusBadRequest:
			msg = "客户端请求错误"
			level = slog.LevelWarn
		case latency > 500*time.Millisecond:
			msg = "请求处理耗时过长"
			level = slog.LevelWarn
			attrs = append(attrs, slog.Bool("slow_request", true))
		}
		logger.LogAttrs(c.Request.Context(), level, msg, attrs...)
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
