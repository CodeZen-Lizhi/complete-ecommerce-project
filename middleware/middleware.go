package middleware

import (
	"bufio"
	"context"
	"ecommerce/container"
	"ecommerce/handler"
	"ecommerce/util"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CorsMiddleware 处理跨域请求
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			// 允许携带cookie（仅在明确 Origin 时才允许）
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			// 非浏览器场景保持兼容
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Role")
		// 允许暴露的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

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
		requestID := c.GetString("X-Request-ID")
		// 获取Authorization头
		// 1. 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		tokenString, err := parseBearerToken(authHeader)
		if err != nil {
			slog.WarnContext(
				c.Request.Context(),
				"鉴权头解析失败",
				slog.String("request_id", requestID),
				slog.String("path", c.Request.URL.Path),
				slog.String("method", c.Request.Method),
				slog.String("error", err.Error()),
			)
			handler.AuthError(c, err.Error())
			return
		}
		// 2. 调用工具类解析令牌
		userID, err := util.ParseToken(tokenString)
		if err != nil {
			slog.WarnContext(
				c.Request.Context(),
				"Token校验失败",
				slog.String("request_id", requestID),
				slog.String("path", c.Request.URL.Path),
				slog.String("method", c.Request.Method),
				slog.String("error", err.Error()),
			)
			handler.AuthError(c, "未授权访问，请先登录")
			return
		}
		// 可以在这里解析token获取用户信息并设置到上下文
		c.Set(util.CurrentUserId, userID)
		c.Next()
	}
}

// parseBearerToken 解析 Authorization 头中的 Bearer Token。
func parseBearerToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", errors.New("缺少Authorization头")
	}
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("Authorization格式错误（应为Bearer <token>）")
	}
	return parts[1], nil
}

// AdminAuthMiddleware 验证是否为管理员
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Role头
		role := c.GetHeader("Role")

		// 验证是否为管理员
		if role != "admin" {
			handler.ForbiddenError(c, "没有管理员权限，无法访问")
			return
		}

		c.Next()
	}
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

// Write 重写字节写入并累计响应字节数。
func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

// WriteString 重写字符串写入并累计响应字节数。
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
				panicErr := toPanicError(r)
				rootFrame, stackLines, stackLocations, rawStack := filterStack(string(debug.Stack()))
				requestID := c.GetString("X-Request-ID")
				logAttrs := []any{
					slog.String("panic_type", fmt.Sprintf("%T", r)),
					slog.String("error", panicErr.Error()),
					slog.String("error_chain", buildErrorChain(panicErr)),
					slog.String("request_id", requestID),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
				}
				if len(stackLines) > 0 {
					logAttrs = append(logAttrs, slog.Any("stacktrace", stackLines))
				}
				if rootFrame != nil {
					logAttrs = append(
						logAttrs,
						slog.String("root_function", shortFunctionName(rootFrame.Function)),
						slog.String("root_file", toRelativeProjectPath(rootFrame.File)),
						slog.Int("root_line", rootFrame.Line),
						slog.String("root_location", formatLocation(rootFrame.File, rootFrame.Line)),
					)
				}
				if len(stackLocations) > 0 {
					logAttrs = append(logAttrs, slog.Any("stacktrace_locations", stackLocations))
				}
				if len(stackLines) == 0 {
					logAttrs = append(logAttrs, slog.String("stacktrace_raw", rawStack))
				}
				logger.ErrorContext(
					c.Request.Context(),
					"Panic Recovered",
					logAttrs...,
				)
				emitClickableLocations(rootFrame, stackLocations)
				sendErrorResponse(c, requestID)
			}
		}()
		c.Next()
	}
}

// filterStack 返回业务根因帧 + 可读调用栈 + 可点击定位 + 原始调用栈
func filterStack(stack string) (*runtime.Frame, []string, []string, string) {
	frames := runtimeCallStackFrames(4)
	projectModulePath := findModulePath()
	projectRoot := findProjectRoot()

	businessFrames := make([]runtime.Frame, 0, len(frames))
	for _, frame := range frames {
		if !isProjectFrame(frame, projectModulePath, projectRoot) {
			continue
		}
		if strings.Contains(filepath.ToSlash(frame.File), "/middleware/") {
			continue
		}
		businessFrames = append(businessFrames, frame)
	}
	// 如果过滤后没有业务帧，降级使用原始堆栈，避免丢失信息
	if len(businessFrames) == 0 {
		return nil, nil, nil, stack
	}
	root := businessFrames[0]
	limit := 20
	if len(businessFrames) < limit {
		limit = len(businessFrames)
	}
	stackLines := make([]string, 0, limit)
	stackLocations := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		frame := businessFrames[i]
		stackLines = append(
			stackLines,
			fmt.Sprintf("at %s (%s:%d)", shortFunctionName(frame.Function), toRelativeProjectPath(frame.File), frame.Line),
		)
		stackLocations = append(stackLocations, formatLocation(frame.File, frame.Line))
	}
	return &root, stackLines, stackLocations, stack
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

// findProjectRoot 获取当前进程工作目录并转为统一路径格式。
func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.ToSlash(wd)
}

// runtimeCallStackFrames 读取当前调用栈帧信息。
func runtimeCallStackFrames(skip int) []runtime.Frame {
	pcs := make([]uintptr, 64)
	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return nil
	}
	frames := runtime.CallersFrames(pcs[:n])
	result := make([]runtime.Frame, 0, n)
	for {
		frame, more := frames.Next()
		result = append(result, frame)
		if !more {
			break
		}
	}
	return result
}

// isProjectFrame 判断栈帧是否属于当前项目代码。
func isProjectFrame(frame runtime.Frame, modulePath, projectRoot string) bool {
	if modulePath != "" && strings.HasPrefix(frame.Function, modulePath+"/") {
		return true
	}
	if projectRoot == "" {
		return false
	}
	file := filepath.ToSlash(frame.File)
	return strings.HasPrefix(file, projectRoot+"/")
}

// toPanicError 将 recover 的任意值转换为 error。
func toPanicError(r any) error {
	if err, ok := r.(error); ok {
		return err
	}
	return fmt.Errorf("%v", r)
}

// buildErrorChain 构建错误包装链文本。
func buildErrorChain(err error) string {
	if err == nil {
		return ""
	}
	var lines []string
	current := err
	for depth := 0; current != nil && depth < 10; depth++ {
		lines = append(lines, current.Error())
		current = errors.Unwrap(current)
	}
	return strings.Join(lines, " <- ")
}

// toRelativeProjectPath 将绝对路径转换为相对项目路径。
func toRelativeProjectPath(path string) string {
	path = filepath.ToSlash(path)
	projectRoot := findProjectRoot()
	if projectRoot == "" {
		return path
	}
	rel, err := filepath.Rel(projectRoot, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

// shortFunctionName 提取函数名的短显示形式。
func shortFunctionName(function string) string {
	function = strings.TrimSpace(function)
	if function == "" {
		return function
	}
	if idx := strings.LastIndex(function, "/"); idx >= 0 && idx+1 < len(function) {
		return function[idx+1:]
	}
	return function
}

// formatLocation 组合文件与行号为可读字符串。
func formatLocation(file string, line int) string {
	return fmt.Sprintf("%s:%d", filepath.ToSlash(file), line)
}

// emitClickableLocations 输出纯文本 file:line，提升 IDE 控制台超链接识别率
func emitClickableLocations(rootFrame *runtime.Frame, stackLocations []string) {
	const maxClickableStackFrames = 5

	rootLocation := ""
	if rootFrame != nil {
		rootLocation = formatLocation(rootFrame.File, rootFrame.Line)
		fmt.Fprintf(os.Stderr, "PANIC_ROOT %s\n", rootLocation)
	}
	if len(stackLocations) == 0 {
		return
	}
	seen := make(map[string]struct{}, len(stackLocations))
	count := 0
	for _, location := range stackLocations {
		if location == rootLocation {
			continue
		}
		if _, ok := seen[location]; ok {
			continue
		}
		seen[location] = struct{}{}
		fmt.Fprintf(os.Stderr, "PANIC_STACK %s\n", location)
		count++
		if count >= maxClickableStackFrames {
			break
		}
	}
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
		latencyMs := latency.Milliseconds()
		statusCode := c.Writer.Status()
		route := c.FullPath()
		if route == "" {
			route = req.URL.Path
		}

		// 使用 slog.Attr 组织日志字段，更具可读性
		attrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.String("route", route),
			slog.Int("status_code", statusCode),
			slog.Duration("latency", latency),
			slog.Int64("latency_ms", latencyMs),
			slog.Int("response_size", w.size),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", req.UserAgent()),
			slog.String("referer", req.Referer()),
		}
		if req.ContentLength >= 0 {
			attrs = append(attrs, slog.Int64("content_length", req.ContentLength))
		}
		if businessCode, ok := c.Get(util.ResponseCodeKey); ok {
			attrs = append(attrs, slog.Any("business_code", businessCode))
		}
		if userID, ok := c.Get(util.CurrentUserId); ok {
			attrs = append(attrs, slog.Any("user_id", userID))
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs,
				slog.String("error", c.Errors.String()),
				slog.Int("error_count", len(c.Errors)),
			)
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

// sendErrorResponse 返回标准化JSON响应（前后端统一格式）
func sendErrorResponse(c *gin.Context, requestID string) {
	// 避免重复响应（若业务Handler已写响应，跳过）
	if c.Writer.Written() {
		return
	}
	handler.SystemError(c, "服务器内部错误，请联系技术支持", requestID)
}
