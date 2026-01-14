package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"smtp-mail/backend/config"
	"smtp-mail/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize 初始化数据库连接
func Initialize() error {
	cfg := config.GetConfig()
	
	// 确保数据库目录存在
	dbDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接SQLite数据库
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层sql.DB以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 自动迁移所有模型
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	DB = db
	log.Printf("数据库初始化成功: %s", cfg.Database.Path)

	return nil
}

// autoMigrate 自动迁移所有模型
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.SMTPConfig{},
		&models.EmailTemplate{},
		&models.EmailHistory{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}