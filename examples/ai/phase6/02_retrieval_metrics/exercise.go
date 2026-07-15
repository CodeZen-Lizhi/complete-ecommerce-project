package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type retrievalCase struct {
	ID          string
	RelevantIDs map[string]struct{}
}

type retrievalRun struct {
	CaseID    string
	RankedIDs []string
}

// recallAtK 计算前 K 个结果中相关文档的覆盖比例。
func recallAtK(relevant map[string]struct{}, ranked []string, k int) (float64, error) {
	return 0, errExerciseIncomplete
}

// reciprocalRank 计算首个相关结果的倒数排名。
func reciprocalRank(relevant map[string]struct{}, ranked []string) (float64, error) {
	return 0, errExerciseIncomplete
}

// summarizeRetrieval 汇总总体与标签分组的 Recall@K、MRR 和 Hit Rate。
func summarizeRetrieval(cases []retrievalCase, runs []retrievalRun, k int) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“检索指标”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：加载 Golden Dataset 和每个查询的原始排名。
	// TODO 2：校验 K、相关集合、重复结果和缺失 Case。
	// TODO 3：实现 recallAtK 与 reciprocalRank，并用手算小样本验证。
	// TODO 4：实现 summarizeRetrieval，计算总体和分组 Hit Rate、Recall、MRR。
	// TODO 5：输出具体失败 Case 和排名证据，禁止只给平均数。
	return errExerciseIncomplete
}
