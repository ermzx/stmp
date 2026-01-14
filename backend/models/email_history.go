package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// EmailStatus 邮件发送状态
type EmailStatus string

const (
	EmailStatusSuccess EmailStatus = "success"
	EmailStatusFailed  EmailStatus = "failed"
)

// StringSlice 用于存储JSON字符串切片
type StringSlice []string

// Scan 实现sql.Scanner接口
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言失败")
	}
	return json.Unmarshal(bytes, s)
}

// Value 实现driver.Valuer接口
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Attachment 附件信息
type Attachment struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}

// AttachmentSlice 用于存储JSON附件切片
type AttachmentSlice []Attachment

// Scan 实现sql.Scanner接口
func (a *AttachmentSlice) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言失败")
	}
	return json.Unmarshal(bytes, a)
}

// Value 实现driver.Valuer接口
func (a AttachmentSlice) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// EmailHistory 邮件发送历史模型
type EmailHistory struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SmtpConfigID uint           `gorm:"not null;index" json:"smtp_config_id"`
	SmtpConfig   SMTPConfig     `gorm:"foreignKey:SmtpConfigID" json:"smtp_config,omitempty"`
	ToEmail      string         `gorm:"type:varchar(255);not null" json:"to_email"`
	CcEmail      StringSlice    `gorm:"type:text" json:"cc_email"`
	BccEmail     StringSlice    `gorm:"type:text" json:"bcc_email"`
	Subject      string         `gorm:"type:varchar(255);not null" json:"subject"`
	Body         string         `gorm:"type:text;not null" json:"body"`
	Attachments  AttachmentSlice `gorm:"type:text" json:"attachments"`
	Status       EmailStatus    `gorm:"type:varchar(20);not null;default:'failed'" json:"status"`
	ErrorMessage string         `gorm:"type:text" json:"error_message"`
	SentAt       time.Time      `json:"sent_at"`
	CreatedAt    time.Time      `json:"created_at"`
}

// TableName 指定表名
func (EmailHistory) TableName() string {
	return "email_histories"
}

// BeforeCreate GORM钩子：创建前设置时间
func (e *EmailHistory) BeforeCreate(tx *gorm.DB) error {
	if e.SentAt.IsZero() {
		e.SentAt = time.Now()
	}
	return nil
}

// IsSuccess 检查邮件是否发送成功
func (e *EmailHistory) IsSuccess() bool {
	return e.Status == EmailStatusSuccess
}

// IsFailed 检查邮件是否发送失败
func (e *EmailHistory) IsFailed() bool {
	return e.Status == EmailStatusFailed
}