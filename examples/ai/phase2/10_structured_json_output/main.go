package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"
)

var (
	errJSONSyntax    = errors.New("JSON 语法错误")
	errJSONStructure = errors.New("JSON 结构错误")
	errBusinessField = errors.New("业务字段错误")
)

const (
	baseURL      = "http://localhost:8084/v1"
	apiKey       = "sk-qWpIk8nVsa8VGJyNHtcAS4VaMhCDJB1z"
	modelName    = "gpt-5.5"
	modelTimeout = 30 * time.Second
)

type structuredOutputMode string

const (
	promptOnlyMode       structuredOutputMode = "prompt_only"
	nativeJSONSchemaMode structuredOutputMode = "native_json_schema"
	// TODO 12：修改 selectedOutputMode，分别运行 Prompt 与原生 JSON Schema 模式。
	selectedOutputMode structuredOutputMode = nativeJSONSchemaMode
)

type structuredAnswer struct {
	Summary    string   `json:"summary" jsonschema:"minLength=1"`
	KeyPoints  []string `json:"key_points" jsonschema:"minItems=1,minLength=1"`
	Confidence float64  `json:"confidence" jsonschema:"minimum=0,maximum=1"`
}

type structuredAnswerPayload struct {
	Summary    string   `json:"summary"`
	KeyPoints  []string `json:"key_points"`
	Confidence *float64 `json:"confidence"`
}

type modelClient interface {
	// Generate 根据消息生成原始文本；实现不得自动清洗或修复 JSON。
	Generate(ctx context.Context, messages []*schema.Message) (string, error)
}

type einoModelClient struct {
	chatModel einomodel.BaseChatModel
}

// Generate 调用 Eino ChatModel，并原样返回模型文本，不清洗或修复 JSON。
func (m einoModelClient) Generate(ctx context.Context, messages []*schema.Message) (string, error) {
	if ctx == nil {
		return "", errors.New("Context 不能为空")
	}
	if m.chatModel == nil {
		return "", errors.New("Eino ChatModel 未配置")
	}
	if len(messages) == 0 {
		return "", errors.New("消息列表不能为空")
	}

	// TODO 3：通过 Model Client 生成原始文本，并用 %w 保留模型错误链。
	response, err := m.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("Eino Generate 调用失败: %w", err)
	}
	if response == nil {
		return "", errors.New("模型响应不能为空")
	}
	return response.Content, nil
}

// main 运行结构化 JSON Prompt、模型调用、解析和字段校验练习。
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), modelTimeout)
	defer cancel()

	chatModel, err := newChatModel(ctx, selectedOutputMode)
	if err != nil {
		fmt.Printf("创建 ChatModel 失败: %v\n", err)
		return
	}
	client := einoModelClient{chatModel: chatModel}

	answer, err := runExercise(ctx, client, "解释 Go interface", selectedOutputMode)
	if err != nil {
		fmt.Printf("结构化输出练习失败: %v\n", err)
		return
	}

	fmt.Printf("摘要: %s\n关键点: %v\n置信度: %.2f\n", answer.Summary, answer.KeyPoints, answer.Confidence)
}

// newChatModel 校验顶部常量配置，并按结构化输出模式创建可复用的 Eino ChatModel。
func newChatModel(ctx context.Context, mode structuredOutputMode) (einomodel.BaseChatModel, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if err := validateModelConfig(baseURL, apiKey, modelName); err != nil {
		return nil, err
	}

	responseFormat, err := buildResponseFormat(mode)
	if err != nil {
		return nil, fmt.Errorf("构造响应格式失败: %w", err)
	}

	// TODO 14：服务端不支持 JSON Schema 时保留创建或调用错误，不静默降级。
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:        baseURL,
		APIKey:         apiKey,
		Model:          modelName,
		Timeout:        modelTimeout,
		ResponseFormat: responseFormat,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 Eino ChatModel 失败: %w", err)
	}
	return chatModel, nil
}

// validateModelConfig 校验模型常量配置，默认占位值不得创建外部客户端。
func validateModelConfig(baseURLValue, apiKeyValue, modelNameValue string) error {
	if strings.TrimSpace(baseURLValue) == "" {
		return errors.New("Base URL 未配置")
	}
	if strings.TrimSpace(apiKeyValue) == "" || apiKeyValue == "replace-with-your-api-key" {
		return errors.New("API Key 未配置")
	}
	if strings.TrimSpace(modelNameValue) == "" {
		return errors.New("Model 未配置")
	}
	return nil
}

// buildResponseFormat 为原生模式生成严格 JSON Schema；Prompt 模式不设置响应格式。
func buildResponseFormat(mode structuredOutputMode) (*openai.ChatCompletionResponseFormat, error) {
	switch mode {
	case promptOnlyMode:
		return nil, nil
	case nativeJSONSchemaMode:
		// TODO 9：从 structuredAnswer 生成禁止额外字段的 JSON Schema。
		reflector := jsonschema.Reflector{
			Anonymous:      true,
			DoNotReference: true,
		}
		answerSchema := reflector.Reflect(structuredAnswer{})
		answerSchema.Version = ""
		// TODO 10：通过 ResponseFormat 传入 json_schema，并设置 Strict: true。
		return &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:        "structured_answer",
				Description: "Go 学习问题的结构化回答",
				Strict:      true,
				JSONSchema:  answerSchema,
			},
		}, nil
	default:
		return nil, fmt.Errorf("不支持的结构化输出模式: %q", mode)
	}
}

// buildStructuredMessages 按输出模式构造 Prompt 约束或原生 JSON Schema 消息。
func buildStructuredMessages(
	ctx context.Context,
	question string,
	mode structuredOutputMode,
) ([]*schema.Message, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if strings.TrimSpace(question) == "" {
		return nil, errors.New("问题不能为空")
	}
	switch mode {
	case promptOnlyMode:
		return []*schema.Message{
			// TODO 1：在 System Message 中禁止 Markdown 代码块和额外文本。
			schema.SystemMessage("只能返回一个 JSON 对象，禁止 Markdown 代码块、解释文字和额外字段。"),
			// TODO 2：在 User Message 中提供固定 JSON 结构和当前问题。
			schema.UserMessage(fmt.Sprintf(`请回答下面的问题：
%s

只能返回以下结构的 JSON：
{
  "summary": "简短总结",
  "key_points": ["关键点1", "关键点2"],
  "confidence": 0.8
}
`, question)),
		}, nil
	case nativeJSONSchemaMode:
		// TODO 11：原生模式不在 Prompt 中重复 JSON 示例，只保留问题语义。
		return []*schema.Message{
			schema.SystemMessage("你是一个 Go 学习助手，请准确回答问题并遵守调用方提供的响应 Schema。"),
			schema.UserMessage(question),
		}, nil
	default:
		return nil, fmt.Errorf("不支持的结构化输出模式: %q", mode)
	}
}

// runExercise 串联消息构造、模型调用和严格 JSON 解析。
func runExercise(
	ctx context.Context,
	client modelClient,
	question string,
	mode structuredOutputMode,
) (structuredAnswer, error) {
	if ctx == nil {
		return structuredAnswer{}, errors.New("Context 不能为空")
	}
	if client == nil {
		return structuredAnswer{}, errors.New("Model Client 不能为空")
	}

	messages, err := buildStructuredMessages(ctx, question, mode)
	if err != nil {
		return structuredAnswer{}, fmt.Errorf("构造结构化消息失败: %w", err)
	}

	raw, err := client.Generate(ctx, messages)
	if err != nil {
		return structuredAnswer{}, fmt.Errorf("模型生成结构化输出失败: %w", err)
	}

	// TODO 4：把原始文本直接交给严格解析器，不静默截取、清洗或修复。
	answer, err := decodeAndValidate([]byte(raw))
	if err != nil {
		return structuredAnswer{}, fmt.Errorf("解析结构化输出失败: %w", err)
	}
	return answer, nil
}

// decodeAndValidate 解析模型 JSON，并校验所有业务字段。
func decodeAndValidate(raw []byte) (structuredAnswer, error) {
	if len(raw) == 0 {
		return structuredAnswer{}, fmt.Errorf("%w: 模型输出不能为空", errJSONSyntax)
	}

	var payload structuredAnswerPayload
	decoder := json.NewDecoder(bytes.NewReader(raw))
	// TODO 5：使用 Decoder 和 DisallowUnknownFields 拒绝未知字段。
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		// TODO 8：区分语法和结构错误，并用 %w 保留底层错误链。
		var syntaxError *json.SyntaxError
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.As(err, &syntaxError) {
			return structuredAnswer{}, fmt.Errorf("%w: %w", errJSONSyntax, err)
		}
		return structuredAnswer{}, fmt.Errorf("%w: %w", errJSONStructure, err)
	}

	// TODO 6：拒绝第二个 JSON 值和尾随非空内容。
	var extra any
	if err := decoder.Decode(&extra); err == nil {
		return structuredAnswer{}, fmt.Errorf("%w: 模型输出包含多个 JSON 值", errJSONStructure)
	} else if !errors.Is(err, io.EOF) {
		return structuredAnswer{}, fmt.Errorf("%w: JSON 包含尾随内容: %w", errJSONStructure, err)
	}

	// TODO 7：校验必填字段、数组元素和 Confidence 范围。
	if strings.TrimSpace(payload.Summary) == "" {
		return structuredAnswer{}, fmt.Errorf("%w: summary 不能为空", errBusinessField)
	}
	if len(payload.KeyPoints) == 0 {
		return structuredAnswer{}, fmt.Errorf("%w: key_points 不能为空", errBusinessField)
	}
	for index, keyPoint := range payload.KeyPoints {
		if strings.TrimSpace(keyPoint) == "" {
			return structuredAnswer{}, fmt.Errorf("%w: key_points[%d] 不能为空", errBusinessField, index)
		}
	}
	if payload.Confidence == nil {
		return structuredAnswer{}, fmt.Errorf("%w: confidence 不能为空", errBusinessField)
	}
	if *payload.Confidence < 0 || *payload.Confidence > 1 {
		return structuredAnswer{}, fmt.Errorf("%w: confidence 必须位于 [0,1]", errBusinessField)
	}

	return structuredAnswer{
		Summary:    payload.Summary,
		KeyPoints:  payload.KeyPoints,
		Confidence: *payload.Confidence,
	}, nil
}
