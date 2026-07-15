package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	baseURL       = "http://localhost:8084/v1"
	apiKey        = "replace-with-your-api-key"
	modelName     = "gpt-5.4-mini"
	streamTimeout = 30 * time.Second
)

var (
	errExerciseIncomplete                         = errors.New("练习尚未完成，请按 TODO 顺序实现")
	_                     einomodel.BaseChatModel = (*openai.ChatModel)(nil)
)

// main 创建 ChatModel，并把流式响应逐块写到标准输出。
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), streamTimeout)
	defer cancel()

	chatModel, err := newChatModel(ctx)
	if err != nil {
		fmt.Printf("创建 ChatModel 失败: %v\n", err)
		return
	}

	messages := []*schema.Message{
		schema.SystemMessage("你是一个 Go 学习助手，请用简洁中文回答。"),
		schema.UserMessage("请解释流式响应中的 EOF 和取消。"),
	}
	if err := streamAnswer(ctx, chatModel, messages, os.Stdout); err != nil {
		fmt.Printf("流式调用失败: %v\n", err)
	}
}

// newChatModel 校验占位配置，并创建 Eino OpenAI-compatible ChatModel。
func newChatModel(ctx context.Context) (einomodel.BaseChatModel, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if strings.TrimSpace(apiKey) == "" || apiKey == "replace-with-your-api-key" {
		return nil, errors.New("API Key 未配置")
	}

	config := &openai.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
		Timeout: streamTimeout,
	}
	chatModel, err := openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("创建 Eino ChatModel 失败: %w", err)
	}
	return chatModel, nil
}

// streamAnswer 消费 Eino StreamReader，区分正常 EOF、取消和中途错误。
func streamAnswer(
	ctx context.Context,
	chatModel einomodel.BaseChatModel,
	messages []*schema.Message,
	writer io.Writer,
) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}
	if chatModel == nil {
		return errors.New("ChatModel 不能为空")
	}
	if len(messages) == 0 {
		return errors.New("消息列表不能为空")
	}
	if writer == nil {
		return errors.New("输出 Writer 不能为空")
	}

	// TODO 1：调用 chatModel.Stream(ctx, messages)，检查创建流时的错误。
	// TODO 2：成功后立即 defer reader.Close()；当前 Eino 版本 Close 没有返回值。
	// TODO 3：循环调用 reader.Recv()。errors.Is(err, io.EOF) 表示正常结束；
	// Context 取消或 DeadlineExceeded 要保留原因返回，其他中途错误使用 %w 包装。
	// TODO 4：拒绝 nil chunk；把每个 chunk.Content 写入 writer，并检查写入错误。
	// TODO 5：至少收到一个非空文本块才算成功；空流返回明确错误。
	return errExerciseIncomplete
}
