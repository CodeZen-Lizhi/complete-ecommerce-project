package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// tokenUsage 统一描述模型调用的 Token 使用量；未返回时 Available 为 false。
type tokenUsage struct {
	PromptTokens     int  `json:"prompt_tokens"`
	CompletionTokens int  `json:"completion_tokens"`
	TotalTokens      int  `json:"total_tokens"`
	Available        bool `json:"available"`
}

// modelCall 是业务层传入 Provider 的已组装消息，不暴露具体厂商请求类型。
type modelCall struct {
	Messages []*schema.Message
}

// modelResult 是非流式模型调用的统一结果。
type modelResult struct {
	Content string
	Usage   tokenUsage
}

// modelChunk 是流式模型调用的单个文本块与可选最终 Token 使用量。
type modelChunk struct {
	Content string
	Usage   tokenUsage
}

// messageStream 抽象流式读取与关闭，避免业务编排层依赖 Eino Reader。
type messageStream interface {
	// Recv 返回下一个模型文本块；正常结束时返回 io.EOF。
	Recv() (modelChunk, error)
	// Close 释放底层模型流和其持有的调用资源。
	Close()
}

// chatProvider 是分类和回答阶段共享的模型边界。
type chatProvider interface {
	// Generate 执行一次非流式模型调用并返回统一结果。
	Generate(context.Context, modelCall) (modelResult, error)
	// Stream 创建可消费的模型流；调用方必须在结束时关闭它。
	Stream(context.Context, modelCall) (messageStream, error)
}

// providerConfig 保存创建 OpenAI-compatible Provider 所需的配置。
type providerConfig struct {
	BaseURL        string
	APIKey         string
	Model          string
	Timeout        time.Duration
	ResponseFormat *openai.ChatCompletionResponseFormat
}

// einoChatProvider 是 Eino ChatModel 的项目适配器，具体转换由练习 TODO 完成。
type einoChatProvider struct {
	chatModel einomodel.BaseChatModel
}

// newOpenAIProvider 校验配置并创建真实的 Eino OpenAI-compatible ChatModel 入口。
func newOpenAIProvider(ctx context.Context, config providerConfig) (chatProvider, error) {
	if ctx == nil {
		return nil, fmt.Errorf("Provider Context 不能为空")
	}
	if err := validateProviderConfig(config); err != nil {
		return nil, err
	}
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:        config.BaseURL,
		APIKey:         config.APIKey,
		Model:          config.Model,
		Timeout:        config.Timeout,
		ResponseFormat: config.ResponseFormat,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 Eino OpenAI-compatible ChatModel 失败: %w", err)
	}
	return &einoChatProvider{chatModel: chatModel}, nil
}

// validateProviderConfig 在 SDK 构造前验证所有 Provider 配置字段。
func validateProviderConfig(config providerConfig) error {
	if strings.TrimSpace(config.BaseURL) == "" {
		return fmt.Errorf("Provider Base URL 未配置")
	}
	if strings.TrimSpace(config.APIKey) == "" {
		return fmt.Errorf("Provider API Key 未配置")
	}
	if strings.TrimSpace(config.Model) == "" {
		return fmt.Errorf("Provider Model 未配置")
	}
	if config.Timeout <= 0 {
		return fmt.Errorf("Provider Timeout 必须大于零")
	}
	return nil
}

// Generate 将项目 modelCall 转成 Eino 消息并归一化非流式结果。
func (provider *einoChatProvider) Generate(ctx context.Context, call modelCall) (modelResult, error) {
	if ctx == nil {
		return modelResult{}, fmt.Errorf("Generate Context 不能为空")
	}
	if provider == nil || provider.chatModel == nil {
		return modelResult{}, fmt.Errorf("Eino ChatModel 未配置")
	}
	if err := validateModelCall(call); err != nil {
		return modelResult{}, err
	}
	result, err := provider.chatModel.Generate(ctx, call.Messages)
	if err != nil {
		return modelResult{}, fmt.Errorf("Eino Generate 调用失败: %w", err)
	}
	//  1：调用 chatModel.Generate，校验非空文本并从 ResponseMeta 归一化可选 Token Usage。
	if result == nil {
		return modelResult{}, fmt.Errorf("模型响应不能为空")
	}
	if strings.TrimSpace(result.Content) == "" {
		return modelResult{}, fmt.Errorf("模型响应文本不能为空")
	}

	modelResultValue := modelResult{Content: result.Content}
	if result.ResponseMeta != nil && result.ResponseMeta.Usage != nil {
		modelResultValue.Usage = tokenUsage{
			PromptTokens:     result.ResponseMeta.Usage.PromptTokens,
			CompletionTokens: result.ResponseMeta.Usage.CompletionTokens,
			TotalTokens:      result.ResponseMeta.Usage.TotalTokens,
			Available:        true,
		}
	}
	return modelResultValue, nil
}

// Stream 将项目 modelCall 转成 Eino 流，并把 Reader 生命周期封装进 messageStream。
func (provider *einoChatProvider) Stream(ctx context.Context, call modelCall) (messageStream, error) {
	if ctx == nil {
		return nil, fmt.Errorf("Stream Context 不能为空")
	}
	if provider == nil || provider.chatModel == nil {
		return nil, fmt.Errorf("Eino ChatModel 未配置")
	}
	if err := validateModelCall(call); err != nil {
		return nil, err
	}
	reader, err := provider.chatModel.Stream(ctx, call.Messages)
	if err != nil {
		return nil, fmt.Errorf("创建 Eino 模型流失败: %w", err)
	}
	if reader == nil {
		return nil, fmt.Errorf("Eino 模型流为空")
	}
	//  2：调用 chatModel.Stream，并以适配器暴露 Recv/Close、文本块和最终 Usage。
	return &einoMessageStream{reader: reader}, nil
}

// einoMessageStream 将 Eino StreamReader 转换为项目内部的 messageStream。
type einoMessageStream struct {
	reader *schema.StreamReader[*schema.Message]
}

// Recv 读取并转换一个 Eino 消息块；EOF 和底层接收错误原样向上游传播。
func (stream *einoMessageStream) Recv() (modelChunk, error) {
	if stream == nil || stream.reader == nil {
		return modelChunk{}, fmt.Errorf("Eino 流未初始化")
	}

	message, err := stream.reader.Recv()
	if err != nil {
		// io.EOF、context.Canceled 等错误原样传出去
		return modelChunk{}, err
	}
	if message == nil {
		return modelChunk{}, fmt.Errorf("模型流消息为空")
	}

	chunk := modelChunk{
		Content: message.Content,
	}

	// Token Usage 通常只在最后一个空内容块中返回
	if message.ResponseMeta != nil && message.ResponseMeta.Usage != nil {
		usage := message.ResponseMeta.Usage
		chunk.Usage = tokenUsage{
			PromptTokens:     usage.PromptTokens,
			CompletionTokens: usage.CompletionTokens,
			TotalTokens:      usage.TotalTokens,
			Available:        true,
		}
	}

	return chunk, nil
}

// Close 释放底层 Eino 流及其持有的资源。
func (stream *einoMessageStream) Close() {
	if stream != nil && stream.reader != nil {
		stream.reader.Close()
	}
}

// validateModelCall 拒绝空消息列表和空消息内容，避免无意义的模型请求。
func validateModelCall(call modelCall) error {
	if len(call.Messages) == 0 {
		return fmt.Errorf("模型消息不能为空")
	}
	for index, message := range call.Messages {
		if message == nil || strings.TrimSpace(message.Content) == "" {
			return fmt.Errorf("模型消息 %d 不能为空", index)
		}
	}
	return nil
}
