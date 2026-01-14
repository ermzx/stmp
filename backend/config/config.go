package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config 应用程序配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Security SecurityConfig `mapstructure:"security"`
	SMTP     SMTPConfig     `mapstructure:"smtp"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	UploadDir    string   `mapstructure:"upload_dir"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWTSecret      string   `mapstructure:"jwt_secret"`
	JWTExpireHours int      `mapstructure:"jwt_expire_hours"`
	CORSEnabled    bool     `mapstructure:"cors_enabled"`
	CORSOrigins    []string `mapstructure:"cors_origins"`
	BcryptCost     int      `mapstructure:"bcrypt_cost"`
}

// SMTPConfig SMTP默认配置
type SMTPConfig struct {
	DefaultHost  string `mapstructure:"default_host"`
	DefaultPort  int    `mapstructure:"default_port"`
	DefaultUseTLS bool   `mapstructure:"default_use_tls"`
}

var appConfig *Config

// GetConfig 获取配置实例
func GetConfig() *Config {
	if appConfig == nil {
		appConfig = loadConfig()
	}
	return appConfig
}

// loadConfig 加载配置文件
func loadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 设置默认值
	viper.SetDefault("server.port", 7700)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.path", "./data/smtp-mail.db")
	viper.SetDefault("upload.max_size", 10485760)
	viper.SetDefault("upload.upload_dir", "./data/uploads")
	viper.SetDefault("security.jwt_expire_hours", 24)
	viper.SetDefault("security.cors_enabled", true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("配置文件未找到，使用默认配置")
		} else {
			log.Fatalf("读取配置文件失败: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 环境变量覆盖配置（优先级高于配置文件）
	if port := os.Getenv("SERVER_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			config.Server.Port = p
			log.Printf("环境变量覆盖: SERVER_PORT=%d", port)
		}
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		config.Server.Mode = mode
		log.Printf("环境变量覆盖: SERVER_MODE=%s", mode)
	}

	fmt.Printf("配置加载成功: 服务器端口=%d, 模式=%s\n", config.Server.Port, config.Server.Mode)

	return &config
}