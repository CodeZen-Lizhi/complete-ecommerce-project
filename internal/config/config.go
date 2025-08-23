package config

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

// Cfg 全局配置实例
var Cfg *Config

// Config 总配置结构体
type Config struct {
	App   AppConfig   `yaml:"app"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"` // Redis 配置节点
}

// RedisConfig Redis 配置结构体（与 YAML 中 redis 节点字段对应）
type RedisConfig struct {
	Addr            string        `yaml:"addr"`           // 地址：ip:port
	Password        string        `yaml:"password"`       // 密码（无则空）
	DB              int           `yaml:"db"`             // 数据库编号
	PoolSize        int           `yaml:"pool_size"`      // 连接池大小
	ConnMaxIdleTime time.Duration `yaml:"idle_timeout"`   // 空闲超时（秒）
	MinIdleConns    int           `yaml:"min_idle_conns"` // 最小空闲连接数
}

// AppConfig 应用基本配置
type AppConfig struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

// MySQLConfig MySQL数据库配置
type MySQLConfig struct {
	Dsn             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	Debug           bool   `yaml:"debug"`
}

// Init 从配置文件加载配置
func Init() error {
	// 1. 从环境变量获取当前环境（默认 dev）
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	// 2. Viper 配置：读取对应环境的配置文件
	viper.SetConfigName(fmt.Sprintf("config.%s", env)) // 如 config.dev
	viper.SetConfigType("yaml")                        // 明确配置格式
	viper.AddConfigPath("./configs")                   // 配置文件目录
	// 3. 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败：%v", err)
		return err
	}
	// 4. 绑定环境变量（可选，优先级：环境变量 > 配置文件）
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP") // 环境变量前缀，如 APP_REDIS_HOST
	err := viper.BindEnv("redis.host", "REDIS_HOST")
	if err != nil {
		return err
	} // 自定义映射

	// 5. 监听配置文件变化（可选，热更新）
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件更新：%s", e.Name)
		// 重新绑定配置结构体（如需热更新生效）
		err := UnmarshalConfig()
		if err != nil {
			return
		}
	})
	return viper.Unmarshal(&Cfg)
}

// UnmarshalConfig 将Viper读取的配置绑定到全局Cfg结构体，并做校验
func UnmarshalConfig() error {
	// 初始化全局配置结构体
	Cfg = &Config{}

	// 调用Viper的Unmarshal方法，将配置绑定到结构体
	// 注意：结构体字段需用 mapstructure 标签（Viper默认使用mapstructure库解析）
	if err := viper.Unmarshal(Cfg); err != nil {
		return fmt.Errorf("配置绑定失败: %w", err)
	}

	// 配置校验（企业项目必做，避免无效配置导致服务异常）
	if err := validateConfig(Cfg); err != nil {
		return fmt.Errorf("配置校验失败: %w", err)
	}

	return nil
}

// validateConfig 配置校验逻辑（根据业务需求自定义）
func validateConfig(cfg *Config) error {
	// 校验应用配置
	if cfg.App.Port <= 0 || cfg.App.Port > 65535 {
		return errors.New("app.port 必须是1-65535之间的整数")
	}
	if cfg.App.Env == "" {
		return errors.New("app.env 不能为空（dev/test/prod）")
	}

	// 校验MySQL配置
	if cfg.MySQL.Dsn == "" {
		return errors.New("mysql.dsn 不能为空")
	}
	if cfg.MySQL.MaxOpenConns <= 0 {
		return errors.New("mysql.max_open_conns 必须大于0")
	}

	// 校验Redis配置
	if cfg.Redis.Addr == "" {
		return errors.New("redis.addr 不能为空（格式：ip:port）")
	}
	if cfg.Redis.PoolSize <= 0 {
		return errors.New("redis.pool_size 必须大于0")
	}

	return nil
}
