package logger

import (
	"ecommerce/internal/config"
	"log/slog"
	"os"
	"sync"
)

// 全局单例实例
var (
	instance *slog.Logger
	once     sync.Once
)

// InitLogConfig 从配置文件初始化日志
func InitLogConfig() error {
	cfg := config.Cfg.Log
	once.Do(func() {
		// 解析日志级别
		level := parseLevel(cfg.Level)
		// 配置日志处理器
		var handler slog.Handler
		switch cfg.Encoding {
		case "json":
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String("time", a.Value.Time().Format("2006-01-02 15:04:05.000"))
				}
				return a
			}, Level: level})
		default: // console 格式
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String("time", a.Value.Time().Format("2006-01-02 15:04:05.000"))
				}
				return a
			}, Level: level})
		}
		// 创建日志实例
		instance = slog.New(handler)
	})
	return nil
}

// GetLogger 获取全局日志实例，未初始化则返回默认实例
func GetLogger() *slog.Logger {
	if instance == nil {
		return slog.Default()
	}
	// 设置全局默认日志实例
	slog.SetDefault(instance)
	return instance
}

// 解析日志级别字符串为slog.Level
func parseLevel(levelStr string) slog.Level {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default: // 默认info级别
		level = slog.LevelInfo
	}
	return level
}
