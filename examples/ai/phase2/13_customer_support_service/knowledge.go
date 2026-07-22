package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// customerIntent 是经结构化分类校验后的有限业务意图。
type customerIntent string

const (
	intentProductAdvice  customerIntent = "product_advice"
	intentDeliveryReturn customerIntent = "delivery_return"
	intentAfterSales     customerIntent = "after_sales"
	intentGeneral        customerIntent = "general"
)

var supportedIntents = []customerIntent{
	intentProductAdvice,
	intentDeliveryReturn,
	intentAfterSales,
	intentGeneral,
}

// businessKnowledge 保存按已验证意图索引的有限业务知识快照。
type businessKnowledge struct {
	Contexts map[customerIntent]string
}

// loadBusinessKnowledge 从练习固定 JSON 文件加载业务知识，并拒绝缺失或空白分支。
func loadBusinessKnowledge(path string) (businessKnowledge, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return businessKnowledge{}, fmt.Errorf("读取业务知识失败: %w", err)
	}
	var raw map[string]string
	if err := json.Unmarshal(content, &raw); err != nil {
		return businessKnowledge{}, fmt.Errorf("解析业务知识失败: %w", err)
	}
	knowledge := businessKnowledge{Contexts: make(map[customerIntent]string, len(raw))}
	for key, value := range raw {
		knowledge.Contexts[customerIntent(key)] = strings.TrimSpace(value)
	}
	if err := validateBusinessKnowledge(knowledge); err != nil {
		return businessKnowledge{}, err
	}
	return knowledge, nil
}

// validateBusinessKnowledge 确保每个受支持意图都有非空、确定性的本地知识片段。
func validateBusinessKnowledge(knowledge businessKnowledge) error {
	if len(knowledge.Contexts) != len(supportedIntents) {
		return fmt.Errorf("业务知识必须恰好包含 %d 个意图分支", len(supportedIntents))
	}
	for _, intent := range supportedIntents {
		if strings.TrimSpace(knowledge.Contexts[intent]) == "" {
			return fmt.Errorf("业务知识缺少意图 %q", intent)
		}
	}
	for intent := range knowledge.Contexts {
		if !isSupportedIntent(intent) {
			return fmt.Errorf("业务知识包含未知意图 %q", intent)
		}
	}
	return nil
}

// contextFor 返回指定已验证意图的知识；不查询数据库、向量库或外部工具。
func (knowledge businessKnowledge) contextFor(intent customerIntent) (string, error) {
	if !isSupportedIntent(intent) {
		return "", fmt.Errorf("不支持的业务意图 %q", intent)
	}
	contextText := strings.TrimSpace(knowledge.Contexts[intent])
	if contextText == "" {
		return "", fmt.Errorf("意图 %q 缺少业务知识", intent)
	}
	return contextText, nil
}

// isSupportedIntent 判断意图是否属于固定的分类枚举。
func isSupportedIntent(intent customerIntent) bool {
	for _, supported := range supportedIntents {
		if intent == supported {
			return true
		}
	}
	return false
}
