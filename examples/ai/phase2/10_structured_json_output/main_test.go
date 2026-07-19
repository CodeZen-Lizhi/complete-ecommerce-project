package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

// TestNewChatModelRejectsNilContext 验证模型创建不会接受空 Context。
func TestNewChatModelRejectsNilContext(t *testing.T) {
	_, err := newChatModel(nil, promptOnlyMode)
	if err == nil || !strings.Contains(err.Error(), "Context") {
		t.Fatalf("期望 Context 配置错误，实际为 %v", err)
	}
}

// TestValidateModelConfig 验证占位值和缺失配置不会进入外部模型客户端创建阶段。
func TestValidateModelConfig(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		apiKey    string
		modelName string
		wantErr   bool
	}{
		{name: "配置完整", baseURL: "http://localhost:8084/v1", apiKey: "test-key", modelName: "test-model"},
		{name: "BaseURL 缺失", apiKey: "test-key", modelName: "test-model", wantErr: true},
		{name: "APIKey 缺失", baseURL: "http://localhost:8084/v1", modelName: "test-model", wantErr: true},
		{name: "APIKey 为占位符", baseURL: "http://localhost:8084/v1", apiKey: "replace-with-your-api-key", modelName: "test-model", wantErr: true},
		{name: "Model 缺失", baseURL: "http://localhost:8084/v1", apiKey: "test-key", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateModelConfig(test.baseURL, test.apiKey, test.modelName)
			if (err != nil) != test.wantErr {
				t.Fatalf("配置校验结果不符合预期: %v", err)
			}
		})
	}
}

// TestBuildResponseFormat 验证 Prompt 模式与原生 JSON Schema 模式使用不同的模型参数。
// TODO 13：在本文件离线覆盖两种模式、未知字段、多个 JSON、缺失字段、空数组和范围错误。
func TestBuildResponseFormat(t *testing.T) {
	t.Run("Prompt 模式不传 ResponseFormat", func(t *testing.T) {
		format, err := buildResponseFormat(promptOnlyMode)
		if err != nil {
			t.Fatalf("构造 Prompt 模式响应格式失败: %v", err)
		}
		if format != nil {
			t.Fatal("Prompt 模式不应设置 ResponseFormat")
		}
	})

	t.Run("原生模式传入严格 JSON Schema", func(t *testing.T) {
		format, err := buildResponseFormat(nativeJSONSchemaMode)
		if err != nil {
			t.Fatalf("构造原生响应格式失败: %v", err)
		}
		if format == nil || format.Type != openai.ChatCompletionResponseFormatTypeJSONSchema {
			t.Fatal("原生模式未设置 JSON Schema ResponseFormat")
		}
		if format.JSONSchema == nil || !format.JSONSchema.Strict || format.JSONSchema.JSONSchema == nil {
			t.Fatal("原生模式必须提供 Strict JSON Schema")
		}

		rawSchema, err := json.Marshal(format.JSONSchema.JSONSchema)
		if err != nil {
			t.Fatalf("序列化 JSON Schema 失败: %v", err)
		}
		for _, field := range []string{"summary", "key_points", "confidence"} {
			if !strings.Contains(string(rawSchema), `"`+field+`"`) {
				t.Fatalf("JSON Schema 缺少字段 %q: %s", field, rawSchema)
			}
		}
		for _, constraint := range []string{
			`"required"`,
			`"additionalProperties":false`,
			`"minLength":1`,
			`"minItems":1`,
			`"minimum":0`,
			`"maximum":1`,
		} {
			if !strings.Contains(string(rawSchema), constraint) {
				t.Fatalf("JSON Schema 缺少约束 %s: %s", constraint, rawSchema)
			}
		}
	})

	t.Run("未知模式明确失败", func(t *testing.T) {
		_, err := buildResponseFormat(structuredOutputMode("unknown"))
		if err == nil {
			t.Fatal("未知结构化输出模式应返回错误")
		}
	})
}

// TestBuildStructuredMessages 验证两种模式只在 Prompt 模式中重复描述 JSON 结构。
func TestBuildStructuredMessages(t *testing.T) {
	promptMessages, err := buildStructuredMessages(context.Background(), "解释 Go interface", promptOnlyMode)
	if err != nil {
		t.Fatalf("构造 Prompt 模式消息失败: %v", err)
	}
	if len(promptMessages) != 2 || !strings.Contains(promptMessages[1].Content, `"summary"`) {
		t.Fatal("Prompt 模式应在 User Message 中包含 JSON 结构示例")
	}

	nativeMessages, err := buildStructuredMessages(context.Background(), "解释 Go interface", nativeJSONSchemaMode)
	if err != nil {
		t.Fatalf("构造原生模式消息失败: %v", err)
	}
	if len(nativeMessages) != 2 || strings.Contains(nativeMessages[1].Content, `"summary"`) {
		t.Fatal("原生模式不应在 User Message 中重复 JSON 结构")
	}
}

// TestDecodeAndValidate 覆盖合法输出以及语法、结构和业务字段错误。
func TestDecodeAndValidate(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{
			name: "合法输出且 confidence 为零",
			raw:  `{"summary":"接口定义行为契约","key_points":["隐式实现"],"confidence":0}`,
		},
		{
			name:    "JSON 语法错误",
			raw:     `{"summary":`,
			wantErr: errJSONSyntax,
		},
		{
			name:    "未知字段",
			raw:     `{"summary":"总结","key_points":["要点"],"confidence":0.8,"extra":true}`,
			wantErr: errJSONStructure,
		},
		{
			name:    "多个 JSON 值",
			raw:     `{"summary":"总结","key_points":["要点"],"confidence":0.8} {"summary":"第二个"}`,
			wantErr: errJSONStructure,
		},
		{
			name:    "summary 为空白",
			raw:     `{"summary":"   ","key_points":["要点"],"confidence":0.8}`,
			wantErr: errBusinessField,
		},
		{
			name:    "key_points 为空",
			raw:     `{"summary":"总结","key_points":[],"confidence":0.8}`,
			wantErr: errBusinessField,
		},
		{
			name:    "key_points 元素为空白",
			raw:     `{"summary":"总结","key_points":["   "],"confidence":0.8}`,
			wantErr: errBusinessField,
		},
		{
			name:    "confidence 缺失",
			raw:     `{"summary":"总结","key_points":["要点"]}`,
			wantErr: errBusinessField,
		},
		{
			name:    "confidence 越界",
			raw:     `{"summary":"总结","key_points":["要点"],"confidence":1.1}`,
			wantErr: errBusinessField,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			answer, err := decodeAndValidate([]byte(test.raw))
			if test.wantErr == nil {
				if err != nil {
					t.Fatalf("解析合法输出失败: %v", err)
				}
				if answer.Confidence != 0 || len(answer.KeyPoints) != 1 {
					t.Fatalf("解析结果不符合预期: %+v", answer)
				}
				return
			}
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("期望错误 %v，实际为 %v", test.wantErr, err)
			}
		})
	}
}
