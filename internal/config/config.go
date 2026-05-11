package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 应用配置结构体
type Config struct {
	Mode   string       `mapstructure:"mode"`
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// Init 初始化配置
func Init() error {
	// 获取环境变量，默认为 dev
	env := os.Getenv("EAO_ENV")
	if env == "" {
		env = "dev"
	}

	// 验证环境变量值，只允许 dev、test、prod
	validEnvs := map[string]bool{
		"dev":  true,
		"test": true,
		"prod": true,
	}
	if !validEnvs[env] {
		return fmt.Errorf("无效的环境变量 EAO_ENV=%s，只允许: dev, test, prod", env)
	}

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %w", err)
	}

	loader := viper.New()
	addConfigPaths(loader, workDir)
	bindEnvVars(loader)

	if err := readRequiredConfig(loader, "app"); err != nil {
		return fmt.Errorf("读取配置文件 app.yml 失败: %w", err)
	}

	if err := mergeRequiredConfig(loader, env); err != nil {
		log.Printf("未找到配置文件 %s.yml，跳过: %v", env, err)
	}

	if merged, err := mergeOptionalConfig(loader, env+".local"); err != nil {
		return fmt.Errorf("读取配置文件 %s.local.yml 失败: %w", env, err)
	} else if merged {
		log.Printf("已合并本地配置文件: %s.local.yml", env)
	}

	// 将配置解析到结构体
	GlobalConfig = &Config{}
	if err := loader.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	log.Printf("当前环境: %s", env)
	log.Printf("配置文件加载成功: app.yml -> %s.yml", env)
	return nil
}

func addConfigPaths(loader *viper.Viper, workDir string) {
	loader.AddConfigPath(filepath.Join(workDir, "internal/config"))
	loader.AddConfigPath("./internal/config")
	loader.AddConfigPath(".")
}

func bindEnvVars(loader *viper.Viper) {
	loader.SetEnvPrefix("EAO")
	loader.AutomaticEnv()

	loader.BindEnv("mode", "EAO_MODE")

	loader.BindEnv("server.port", "EAO_SERVER_PORT")

	loader.BindEnv("logger.level", "EAO_LOGGER_LEVEL")
	loader.BindEnv("logger.filename", "EAO_LOGGER_FILENAME")
	loader.BindEnv("logger.max_size", "EAO_LOGGER_MAX_SIZE")
	loader.BindEnv("logger.max_age", "EAO_LOGGER_MAX_AGE")
	loader.BindEnv("logger.max_backups", "EAO_LOGGER_MAX_BACKUPS")
}

func readRequiredConfig(loader *viper.Viper, name string) error {
	loader.SetConfigName(name)
	loader.SetConfigType("yml")
	return loader.ReadInConfig()
}

func mergeRequiredConfig(loader *viper.Viper, name string) error {
	loader.SetConfigName(name)
	loader.SetConfigType("yml")
	return loader.MergeInConfig()
}

func mergeOptionalConfig(loader *viper.Viper, name string) (bool, error) {
	loader.SetConfigName(name)
	loader.SetConfigType("yml")
	if err := loader.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return GlobalConfig
}

// GetServerAddr 获取服务器地址
func GetServerAddr() string {
	if GlobalConfig == nil {
		return ":6304" // 默认端口
	}
	return fmt.Sprintf(":%d", GlobalConfig.Server.Port)
}

// GetEnv 获取当前环境（dev/test/prod）
func GetEnv() string {
	env := os.Getenv("EAO_ENV")
	if env == "" {
		return "dev"
	}
	return env
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	return GetEnv() == "prod"
}

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	return GetEnv() == "dev"
}

// IsTest 判断是否为测试环境
func IsTest() bool {
	return GetEnv() == "test"
}
