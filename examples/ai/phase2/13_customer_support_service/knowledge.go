package main

import (
	"fmt"
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

// defaultBusinessKnowledge 返回练习内置的确定性业务知识，不依赖文件或外部服务。
func defaultBusinessKnowledge() businessKnowledge {
	return businessKnowledge{Contexts: map[customerIntent]string{
		intentProductAdvice:  "只根据用户提供的需求给出商品选择维度，例如使用场景、预算、兼容性和售后政策。没有商品目录或实时库存时，明确说明无法确认具体库存、价格或促销。",
		intentDeliveryReturn: "配送时效、运费、退货窗口和退款状态以订单页或官方售后规则为准。不得承诺未核验的到货日期；需要订单信息时引导用户在受保护的订单渠道查询。",
		intentAfterSales:     "商品故障、缺件、破损或超过正常处理时效时，先收集必要订单信息并引导至人工售后。不得要求用户在聊天中发送完整支付凭证、密码或身份证件。",
		intentGeneral:        "回答保持简洁、礼貌和可执行。无法根据有限上下文确认的事实必须明确说明，并在需要订单、账户或支付数据时建议使用受保护的官方渠道。",
	}}
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
