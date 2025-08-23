package mysql

import (
	"fmt"
	"time"

	"ecommerce/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局DB实例
var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() error {
	cfg := config.Cfg.MySQL
	// 构建DSN
	dsn := cfg.Dsn

	// 配置GORM
	gormConfig := &gorm.Config{
		// 根据配置决定是否开启调试模式
		Logger: logger.Default.LogMode(func() logger.LogLevel {
			if cfg.Debug {
				return logger.Info
			}
			return logger.Warn
		}()),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 获取底层sql.DB对象，配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取sql.DB失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	DB = db
	return nil
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
