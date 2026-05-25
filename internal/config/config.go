package config

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Cfg 全局配置实例
var Cfg *Config

// Config 总配置结构体
type Config struct {
	App   AppConfig   `mapstructure:"app"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Log   LogConfig   `mapstructure:"log"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
}

// RedisConfig Redis 配置结构体（与 YAML 中 redis 节点字段对应）
type RedisConfig struct {
	Addr            string        `mapstructure:"addr"`           // 地址：ip:port
	Password        string        `mapstructure:"password"`       // 密码（无则空）
	DB              int           `mapstructure:"db"`             // 数据库编号
	PoolSize        int           `mapstructure:"pool_size"`      // 连接池大小
	ConnMaxIdleTime time.Duration `mapstructure:"idle_timeout"`   // 空闲超时（秒）
	MinIdleConns    int           `mapstructure:"min_idle_conns"` // 最小空闲连接数
}

// AppConfig 应用基本配置
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Port      int    `mapstructure:"port"`
	Env       string `mapstructure:"env"`
	MachineId int64  `mapstructure:"machine_id"`
}

// MySQLConfig MySQL数据库配置
type MySQLConfig struct {
	Dsn             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	Debug           bool   `mapstructure:"debug"`
}

// LogConfig 日志配置。
type LogConfig struct {
	Level     string `mapstructure:"level"`
	Encoding  string `mapstructure:"encoding"`
	AddSource bool   `mapstructure:"add_source"`
}

// JwtConfig JWT配置。
type JwtConfig struct {
	Secret       string `mapstructure:"secret"`
	ExpiresHours int    `mapstructure:"expires_hours"`
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
		slog.Error("读取配置失败", "error", err)
		return err
	}
	//读取普通的配置
	if viper.IsSet("userid") {
		userid := viper.GetString("userid")
		fmt.Println("userid==", userid)
	}
	// 4. 绑定环境变量（可选，优先级：环境变量 > 配置文件）
	viper.SetEnvPrefix("APP") // 环境变量前缀，如 APP_REDIS_HOST
	viper.AutomaticEnv()      //启用Viper自动读取环境变量的功能
	//- 将配置文件中的"redis.addr"字段与"REDIS_HOST"环境变量进行绑定 这样改环境变量就可以覆盖配置文件的数据
	err := viper.BindEnv("redis.addr", "REDIS_HOST")
	if err != nil {
		return err
	} // 自定义映射

	// 5. 初次加载配置并校验
	if err := UnmarshalConfig(); err != nil {
		return err
	}

	// 5. 监听配置文件变化（可选，热更新）
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件更新：%s", e.Name)
		// 重新绑定配置结构体（如需热更新生效）
		if err := UnmarshalConfig(); err != nil {
			log.Printf("配置热更新失败：%v", err)
		}
	})
	return nil
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
	// 校验JWT配置
	if cfg.Jwt.Secret == "" {
		return errors.New("jwt.secret 不能为空")
	}
	if cfg.Jwt.ExpiresHours <= 0 {
		return errors.New("jwt.expires_hours 必须大于0")
	}

	return nil
}
