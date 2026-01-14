package models

import (
	"time"

	"gorm.io/gorm"
)

// EncryptionType 加密类型
type EncryptionType string

const (
	EncryptionNone    EncryptionType = "none"
	EncryptionTLS     EncryptionType = "tls"
	EncryptionStartTLS EncryptionType = "starttls"
)

// SMTPConfig SMTP配置模型
type SMTPConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Host      string         `gorm:"type:varchar(255);not null" json:"host"`
	Port      int            `gorm:"not null" json:"port"`
	Username  string         `gorm:"type:varchar(255)" json:"username"`
	Password  string         `gorm:"type:varchar(255)" json:"password"` // 接收时使用，响应时由代码清除
	FromEmail string         `gorm:"type:varchar(255);not null" json:"from_email"`
	FromName  string         `gorm:"type:varchar(100)" json:"from_name"`
	Encryption EncryptionType `gorm:"type:varchar(20);default:'none'" json:"encryption"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (SMTPConfig) TableName() string {
	return "smtp_configs"
}

// BeforeCreate GORM钩子：创建前确保只有一个默认配置
func (s *SMTPConfig) BeforeCreate(tx *gorm.DB) error {
	if s.IsDefault {
		// 将其他配置的IsDefault设置为false
		tx.Model(&SMTPConfig{}).Where("is_default = ?", true).Update("is_default", false)
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前确保只有一个默认配置
func (s *SMTPConfig) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("IsDefault") && s.IsDefault {
		// 将其他配置的IsDefault设置为false
		tx.Model(&SMTPConfig{}).Where("is_default = ? AND id != ?", true, s.ID).Update("is_default", false)
	}
	return nil
}

// GetDefaultConfig 获取默认SMTP配置
func GetDefaultConfig(db *gorm.DB) (*SMTPConfig, error) {
	var config SMTPConfig
	err := db.Where("is_default = ?", true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}