package models

import (
	"time"

	"gorm.io/gorm"
)

// EmailTemplate 邮件模板模型
type EmailTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Subject   string    `gorm:"type:varchar(255);not null" json:"subject"`
	Body      string    `gorm:"type:text;not null" json:"body"` // 支持HTML内容
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (EmailTemplate) TableName() string {
	return "email_templates"
}

// BeforeCreate GORM钩子：创建前验证
func (e *EmailTemplate) BeforeCreate(tx *gorm.DB) error {
	return e.validate()
}

// BeforeUpdate GORM钩子：更新前验证
func (e *EmailTemplate) BeforeUpdate(tx *gorm.DB) error {
	return e.validate()
}

// validate 验证模板数据
func (e *EmailTemplate) validate() error {
	if e.Name == "" {
		return ErrTemplateNameRequired
	}
	if e.Subject == "" {
		return ErrTemplateSubjectRequired
	}
	if e.Body == "" {
		return ErrTemplateBodyRequired
	}
	return nil
}

// 错误定义
var (
	ErrTemplateNameRequired    = &ValidationError{Field: "name", Message: "模板名称不能为空"}
	ErrTemplateSubjectRequired = &ValidationError{Field: "subject", Message: "模板主题不能为空"}
	ErrTemplateBodyRequired    = &ValidationError{Field: "body", Message: "模板内容不能为空"}
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}