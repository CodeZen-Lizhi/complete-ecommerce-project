package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxChatRequestBytes = int64(maxMessageLength + 1024)

// newCustomerSupportRouter 创建独立练习的唯一 HTTP 路由，不接入生产 Router。
func newCustomerSupportRouter(service *customerSupportService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), demoIdentityMiddleware())
	router.POST("/api/ai/chat/stream", handleChatStream(service))
	return router
}

// handleChatStream 在开始 SSE 前完成身份和 JSON 校验，之后委托服务层编排。
func handleChatStream(service *customerSupportService) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		if service == nil {
			writeJSONError(ginContext, http.StatusServiceUnavailable, "service_unavailable", "客服服务暂不可用")
			return
		}
		userID, err := demoUserIDFromContext(ginContext.Request.Context())
		if err != nil {
			writeJSONError(ginContext, http.StatusUnauthorized, "unauthenticated", "缺少或无效的开发身份")
			return
		}
		if !isJSONContentType(ginContext.GetHeader("Content-Type")) {
			writeJSONError(ginContext, http.StatusUnsupportedMediaType, "unsupported_media_type", "请求必须使用 application/json")
			return
		}
		ginContext.Request.Body = http.MaxBytesReader(ginContext.Writer, ginContext.Request.Body, maxChatRequestBytes)
		request, err := decodeChatRequest(ginContext.Request.Body)
		if err != nil {
			writeJSONError(ginContext, http.StatusBadRequest, "invalid_request", "请求体无效")
			return
		}
		if err := validateChatRequest(request); err != nil {
			writeJSONError(ginContext, http.StatusBadRequest, "invalid_request", "会话或消息无效")
			return
		}

		emitter := newSSEEmitter(ginContext)
		if err := service.Stream(ginContext.Request.Context(), userID, request, emitter); err != nil {
			if emitter.Started() {
				_ = emitter.Send("error", publicStreamError(err))
				return
			}
			status, code, message := publicHTTPError(err)
			writeJSONError(ginContext, status, code, message)
		}
	}
}

// isJSONContentType 仅接受 application/json 及其合法参数形式。
func isJSONContentType(value string) bool {
	mediaType, _, err := mime.ParseMediaType(value)
	return err == nil && mediaType == "application/json"
}

// decodeChatRequest 严格解码请求体，拒绝伪造 user_id、未知字段和多个 JSON 值。
func decodeChatRequest(reader io.Reader) (chatRequest, error) {
	if reader == nil {
		return chatRequest{}, fmt.Errorf("请求体不能为空")
	}
	var payload struct {
		SessionID string `json:"session_id"`
		Message   string `json:"message"`
	}
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		return chatRequest{}, fmt.Errorf("解码请求 JSON 失败: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return chatRequest{}, fmt.Errorf("请求体只能包含一个 JSON 值")
		}
		return chatRequest{}, fmt.Errorf("请求体包含尾随内容: %w", err)
	}
	return chatRequest{SessionID: strings.TrimSpace(payload.SessionID), Message: payload.Message}, nil
}

// apiError 是 SSE 开始前返回的稳定 JSON 错误结构。
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// streamError 是 SSE 已开始后的稳定错误事件负载。
type streamError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSONError 写入安全错误响应并终止 Gin 链路。
func writeJSONError(ginContext *gin.Context, status int, code string, message string) {
	ginContext.AbortWithStatusJSON(status, apiError{Code: code, Message: message})
}

// publicHTTPError 将内部错误映射为 SSE 尚未开始时可安全公开的响应。
func publicHTTPError(err error) (int, string, string) {
	if errors.Is(err, errExerciseIncomplete) {
		return http.StatusServiceUnavailable, "exercise_incomplete", "练习核心步骤尚未完成"
	}
	return http.StatusBadGateway, "upstream_failure", "客服服务暂时无法完成请求"
}

// publicStreamError 将 SSE 开始后的内部错误映射为安全事件，不泄露 Redis 或模型细节。
func publicStreamError(err error) streamError {
	_, code, message := publicHTTPError(err)
	return streamError{Code: code, Message: message}
}

// sseEmitter 延迟提交 HTTP 响应，确保开始流前的失败仍可返回普通 JSON 错误。
type sseEmitter struct {
	ginContext *gin.Context
	started    bool
}

// newSSEEmitter 创建尚未写入响应头的 SSE 事件发送器。
func newSSEEmitter(ginContext *gin.Context) *sseEmitter {
	return &sseEmitter{ginContext: ginContext}
}

// Started 报告是否已写出任一 SSE 事件。
func (emitter *sseEmitter) Started() bool {
	return emitter != nil && emitter.started
}

// Send 编码并发送一个固定名称的 SSE 事件，首次发送时才提交响应头。
func (emitter *sseEmitter) Send(event string, payload any) error {
	if emitter == nil || emitter.ginContext == nil {
		return fmt.Errorf("SSE 发送器未初始化")
	}
	if !isAllowedSSEEvent(event) {
		return fmt.Errorf("不支持的 SSE 事件 %q", event)
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("编码 SSE 数据失败: %w", err)
	}
	if !emitter.started {
		emitter.ginContext.Header("Content-Type", "text/event-stream")
		emitter.ginContext.Header("Cache-Control", "no-cache")
		emitter.ginContext.Header("Connection", "keep-alive")
		emitter.ginContext.Status(http.StatusOK)
		emitter.ginContext.Writer.WriteHeaderNow()
		emitter.started = true
	}
	if _, err := fmt.Fprintf(emitter.ginContext.Writer, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return fmt.Errorf("写入 SSE 事件失败: %w", err)
	}
	emitter.ginContext.Writer.Flush()
	return nil
}

// isAllowedSSEEvent 将服务层事件限定为协议约定的四种类型。
func isAllowedSSEEvent(event string) bool {
	switch event {
	case "meta", "delta", "done", "error":
		return true
	default:
		return false
	}
}
