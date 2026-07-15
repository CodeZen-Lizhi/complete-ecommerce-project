package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL = "http://localhost:8084/v1"
	apiKey  = "replace-with-your-api-key"
	model   = "gpt-5.4-mini"

	maxResponseBodyBytes = 4 << 20
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startedAt := time.Now()
	result, statusCode, err := callChat(ctx, []message{
		{Role: "system", Content: "你是一个 Go 学习助手，请用简洁中文回答。"},
		{Role: "user", Content: "请用一个例子解释 Go 的 context.Context 有什么作用。"},
	})
	if err != nil {
		fmt.Printf("调用失败: %v\n", err)
		return
	}

	fmt.Printf("HTTP 状态码: %d\n", statusCode)
	fmt.Printf("模型回答: %s\n", result.Choices[0].Message.Content)
	fmt.Printf("输入 Token: %d\n", result.Usage.PromptTokens)
	fmt.Printf("输出 Token: %d\n", result.Usage.CompletionTokens)
	fmt.Printf("总 Token: %d\n", result.Usage.TotalTokens)
	fmt.Printf("总耗时: %s\n", time.Since(startedAt))
}

func callChat(ctx context.Context, messages []message) (chatResponse, int, error) {
	// TODO 1：检查 apiKey 是否还是占位符或空字符串。
	// 如果没有配置真实 Key，返回空 chatResponse、状态码 0 和一个说明配置错误的 error。
	if strings.TrimSpace(apiKey) == "" || apiKey == "replace-with-your-api-key" {
		return chatResponse{}, 0, fmt.Errorf("API Key 未配置")
	}
	// TODO 2：使用 model、messages 和 Temperature 组装一个 chatRequest 结构体。
	// Temperature 可以先设置为 0.2，让练习结果相对稳定。
	request := chatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.2,
	}
	// TODO 3：使用 encoding/json 包的 Marshal 函数，把 chatRequest 转换为 JSON 字节切片。
	// 必须检查 Marshal 返回的 err；失败时返回空结果，并用 fmt.Errorf 和 %w 包装原始错误。
	req, err := json.Marshal(request)
	if err != nil {
		return chatResponse{}, 0, fmt.Errorf("marshal json: %w", err)
	}
	// TODO 4：生成完整请求地址 endpoint。
	// 先使用 strings.TrimRight 去掉 baseURL 末尾可能存在的“/”，再拼接“/chat/completions”。
	url := strings.TrimRight(baseURL, "/") + "/chat/completions"
	// TODO 5：创建一个携带 ctx 的 HTTP POST 请求。
	// 使用 http.NewRequestWithContext；请求体需要用 bytes.NewReader 把第 3 步的 JSON 字节转换为 io.Reader。
	// 必须检查创建请求时返回的 err。
	requestWithContext, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(req))
	if err != nil {
		return chatResponse{}, 0, fmt.Errorf("create request: %w", err)
	}
	// TODO 6：为请求设置两个 Header。
	// Authorization 的值是“Bearer ”加 apiKey；Content-Type 的值是“application/json”。
	requestWithContext.Header.Set("Content-Type", "application/json")
	requestWithContext.Header.Set("Authorization", "Bearer "+apiKey)

	// TODO 7：使用 http.DefaultClient.Do 发送请求，并接收 resp 和 err。
	// 网络失败、超时或 ctx 被取消时会返回 err，此时还没有可靠的 HTTP 状态码，所以状态码返回 0。
	resp, err := http.DefaultClient.Do(requestWithContext)
	if err != nil {
		return chatResponse{}, 0, fmt.Errorf("do request: %w", err)
	}
	// TODO 8：成功获得 resp 后，立即使用 defer 注册 resp.Body.Close()。
	// 这样无论函数之后从哪个分支返回，响应体最终都会被关闭。
	defer resp.Body.Close()

	// TODO 9：读取响应体。
	// 使用 io.LimitReader 最多读取 4 MiB + 1 字节，再使用 io.ReadAll 得到 responseBody。
	// 多出的 1 字节用于判断响应是否超过上限；读取失败时要返回状态码和包装后的错误。
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes+1))
	if err != nil {
		return chatResponse{}, resp.StatusCode, fmt.Errorf("读取响应失败: %w", err)
	}
	if len(responseBody) > maxResponseBodyBytes {
		return chatResponse{}, resp.StatusCode, fmt.Errorf("模型响应超过 4 MiB 限制")
	}
	// TODO 10：检查 HTTP 状态码是否在 [200, 300) 范围内。
	// 如果不是成功状态码，返回空结果、resp.StatusCode，以及包含响应内容的错误，方便观察服务端报错。
	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return chatResponse{}, resp.StatusCode,
			fmt.Errorf(
				"模型返回错误，状态码 %d: %s",
				resp.StatusCode,
				string(responseBody),
			)
	}
	// TODO 11：声明一个 chatResponse 变量 result，再使用 json.Unmarshal 把 responseBody 解析到 &result。
	// 注意必须传指针；解析失败时返回空结果、resp.StatusCode 和包装后的错误。
	var result chatResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return chatResponse{}, resp.StatusCode, fmt.Errorf("unmarshal response: %w", err)
	}
	// TODO 12：检查 result.Choices 是否为空。
	// 如果没有任何候选回答，不能访问 result.Choices[0]，应返回一个明确错误，避免数组越界 panic。
	if len(result.Choices) == 0 {
		return chatResponse{}, resp.StatusCode,
			fmt.Errorf("模型响应中没有 choices")
	}
	// TODO 13：所有步骤都成功后，返回 result、resp.StatusCode 和 nil。
	return result, resp.StatusCode, nil
}
