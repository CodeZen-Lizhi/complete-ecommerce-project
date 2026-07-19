package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/model/agenticopenai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	baseURL        = "http://localhost:8084/v1"
	modelName      = "gpt-5.4-mini"
	apiKey         = "replace-with-your-api-key"
	requestTimeout = 30 * time.Second
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type comparisonScenario struct {
	Name     string
	Messages []*schema.AgenticMessage
}

var _ einomodel.AgenticModel = (*agenticopenai.ResponsesModel)(nil)

// main 运行 Zero-shot 与 Few-shot 的真实模型对比练习。
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if err := runExercise(ctx); err != nil {
		fmt.Printf("Few-shot 对比练习失败: %v\n", err)
	}
}

// buildComparisonScenarios 固定模型问题，只改变是否提供示例消息。
func buildComparisonScenarios() ([]comparisonScenario, error) {
	return nil, errExerciseIncomplete
}

// newAgenticModel 从顶部常量读取配置并创建真实 Eino ResponsesModel。
func newAgenticModel(ctx context.Context) (einomodel.AgenticModel, error) {
	return nil, errExerciseIncomplete
}

// extractAssistantText 从真实模型响应中提取非空 AssistantGenText。
func extractAssistantText(message *schema.AgenticMessage) (string, error) {
	return "", errExerciseIncomplete
}

// runExercise 按相同问题依次执行 Zero-shot 与 Few-shot 场景并输出对比结果。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义一个输出格式容易偏离的分类问题，以及两组 User/Assistant 示例。
	// TODO 2：实现 buildComparisonScenarios；Zero-shot 只含 System/User，Few-shot 按 System/User/Assistant/User/Assistant/User 排列。
	// TODO 3：校验顶部 apiKey 占位符，并用 baseURL、modelName、requestTimeout 创建真实 ResponsesModel。
	// TODO 4：使用同一个 AgenticModel 依次调用两个场景，除示例消息外不得改变问题、模型和参数。
	// TODO 5：实现 extractAssistantText，拒绝 nil 响应、空文本和错误内容块类型。
	// TODO 6：输出两组原始回答，对比格式遵循率、标签稳定性、示例照抄和错误泛化。
	return errExerciseIncomplete
}
