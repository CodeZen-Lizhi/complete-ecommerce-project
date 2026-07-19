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

type comparisonFixture struct {
	Question string
	Examples [][2]string
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

// newComparisonFixture 定义固定问题和 Few-shot 示例对。
func newComparisonFixture() (comparisonFixture, error) {
	// TODO 1：定义一个输出格式容易偏离的分类问题，以及两组 User/Assistant 示例。
	return comparisonFixture{}, errExerciseIncomplete
}

// buildComparisonScenarios 固定模型问题，只改变是否提供示例消息。
func buildComparisonScenarios(fixture comparisonFixture) ([]comparisonScenario, error) {
	// TODO 2：Zero-shot 只含 System/User，Few-shot 按 System/User/Assistant/User/Assistant/User 排列。
	return nil, errExerciseIncomplete
}

// newAgenticModel 从顶部常量读取配置并创建真实 Eino ResponsesModel。
func newAgenticModel(ctx context.Context) (einomodel.AgenticModel, error) {
	// TODO 3：校验 apiKey 占位符，并用 baseURL、modelName、requestTimeout 创建 ResponsesModel。
	return nil, errExerciseIncomplete
}

// runComparison 使用同一个模型依次调用两个场景。
func runComparison(ctx context.Context, agenticModel einomodel.AgenticModel, scenarios []comparisonScenario) ([]*schema.AgenticMessage, error) {
	// TODO 4：除示例消息外保持问题、模型和参数一致，并保留场景名称对应的错误。
	return nil, errExerciseIncomplete
}

// extractAssistantText 从真实模型响应中提取非空 AssistantGenText。
func extractAssistantText(message *schema.AgenticMessage) (string, error) {
	// TODO 5：拒绝 nil 响应、空文本和错误内容块类型。
	return "", errExerciseIncomplete
}

// printComparison 输出原始回答和对比观察项。
func printComparison(scenarios []comparisonScenario, responses []*schema.AgenticMessage) error {
	// TODO 6：对比格式遵循率、标签稳定性、示例照抄和错误泛化。
	return errExerciseIncomplete
}

// runExercise 按相同问题依次执行 Zero-shot 与 Few-shot 场景并输出对比结果。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	fixture, err := newComparisonFixture()
	if err != nil {
		return fmt.Errorf("准备对比数据失败: %w", err)
	}
	scenarios, err := buildComparisonScenarios(fixture)
	if err != nil {
		return fmt.Errorf("构造对比场景失败: %w", err)
	}
	agenticModel, err := newAgenticModel(ctx)
	if err != nil {
		return fmt.Errorf("创建模型失败: %w", err)
	}
	responses, err := runComparison(ctx, agenticModel, scenarios)
	if err != nil {
		return fmt.Errorf("执行模型对比失败: %w", err)
	}
	return printComparison(scenarios, responses)
}
