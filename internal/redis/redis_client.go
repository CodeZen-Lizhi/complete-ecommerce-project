package redis

import (
	"context"
	"ecommerce/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

// 全局Redis客户端实例
var (
	client *redis.Client
	once   sync.Once // 确保客户端只初始化一次
)

// Init 初始化Redis客户端
// 建议在程序启动时调用（如main函数中）
func Init() error {
	var err error
	once.Do(func() {
		// 从配置文件读取Redis配置
		var cfg = config.Cfg.Redis

		// 创建客户端实例
		client = redis.NewClient(&redis.Options{
			Addr:            cfg.Addr,
			Password:        cfg.Password,
			DB:              cfg.DB,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, pingErr := client.Ping(ctx).Result(); pingErr != nil {
			err = fmt.Errorf("连接Redis失败: %v", pingErr)
			client = nil // 连接失败时置空
		}
	})
	return err
}

// Client 获取全局Redis客户端实例
// 必须在Init()调用成功后使用
func Client() *redis.Client {
	if client == nil {
		panic("Redis客户端未初始化，请先调用Init()")
	}
	return client
}

// Close 关闭Redis连接（程序退出时调用）
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
