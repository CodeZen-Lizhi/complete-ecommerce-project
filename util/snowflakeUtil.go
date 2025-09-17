package util

import (
	"ecommerce/internal/logger"
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

func Init(nodeId int64) {
	//设置自定义起始时间
	customEpoch := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	snowflake.Epoch = customEpoch
	// 创建节点（支持通过选项配置更多参数
	var err error
	node, err = snowflake.NewNode(nodeId)
	if err != nil {
		// 记录错误日志
		logger.GetLogger().Error("创建Snowflake节点失败", "nodeId", nodeId, "error", err)
		return
	}
}

func GenID() int64 {
	if node == nil {
		return 0
	}
	generate := node.Generate()
	return int64(generate)
}

// ParseID 解析雪花ID，返回生成时间等信息
func ParseID(id int64) map[string]interface{} {
	snowflakeID := snowflake.ID(id)
	return map[string]interface{}{
		"timestamp": time.UnixMilli(snowflakeID.Time()).Format(time.DateTime),
		"node":      snowflakeID.Node(),
	}
}
