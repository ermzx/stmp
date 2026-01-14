package services

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"smtp-mail/backend/database"
	"smtp-mail/backend/models"
	"smtp-mail/backend/utils"
)

// EmailService 邮件服务
type EmailService struct {
	smtpService *SMTPService
}

// NewEmailService 创建邮件服务实例
func NewEmailService() *EmailService {
	return &EmailService{
		smtpService: NewSMTPService(),
	}
}

// SendEmailRequest 发送邮件请求
type SendEmailRequest struct {
	SmtpConfigID uint         `json:"smtp_config_id" binding:"required"`
	To           []string     `json:"to" binding:"required,min=1"`
	Cc           []string     `json:"cc"`
	Bcc          []string     `json:"bcc"`
	Subject      string       `json:"subject" binding:"required"`
	Body         string       `json:"body" binding:"required"`
	Attachments  []Attachment `json:"attachments"`
}

// Attachment 附件（用于请求）
type Attachment struct {
	Filename    string `json:"filename" binding:"required"`
	Content     string `json:"content" binding:"required"` // base64编码的文件内容
	ContentType string `json:"content_type"`               // 内容类型（可选）
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(req *SendEmailRequest) (*models.EmailHistory, error) {
	utils.Infof("开始发送邮件: SmtpConfigID=%d, To=%v, Subject=%s, Attachments=%d",
		req.SmtpConfigID, req.To, req.Subject, len(req.Attachments))
	
	// 1. 获取SMTP配置（包含密码）
	config, err := s.smtpService.GetConfigByIDWithPassword(req.SmtpConfigID)
	if err != nil {
		utils.Errorf("获取SMTP配置失败 (ID: %d): %v", req.SmtpConfigID, err)
		return nil, fmt.Errorf("获取SMTP配置失败: %w", err)
	}

	utils.Infof("获取SMTP配置成功: Host=%s, Port=%d, FromEmail=%s", config.Host, config.Port, config.FromEmail)

	// 2. 解密SMTP密码
	password, err := s.smtpService.cryptoService.DecryptPassword(config.Password)
	if err != nil {
		utils.Errorf("解密密码失败: %v", err)
		return nil, fmt.Errorf("解密密码失败: %w", err)
	}

	// 3. 验证收件人邮箱格式
	if err := validateEmails(req.To); err != nil {
		return nil, fmt.Errorf("收件人邮箱格式错误: %w", err)
	}
	if err := validateEmails(req.Cc); err != nil {
		return nil, fmt.Errorf("抄送邮箱格式错误: %w", err)
	}
	if err := validateEmails(req.Bcc); err != nil {
		return nil, fmt.Errorf("密送邮箱格式错误: %w", err)
	}

	// 4. 构建邮件消息
	message, err := s.buildEmailMessage(config, req)
	if err != nil {
		utils.Errorf("构建邮件消息失败: %v", err)
		return nil, fmt.Errorf("构建邮件消息失败: %w", err)
	}

	utils.Infof("邮件消息构建成功: 消息大小=%d 字节", len(message))

	// 5. 发送邮件
	err = s.sendEmailViaSMTP(config, password, req.To, req.Cc, req.Bcc, message)
	if err != nil {
		utils.Errorf("发送邮件失败: %v", err)
		// 记录失败历史
		history := s.createEmailHistory(req, models.EmailStatusFailed, err.Error())
		return history, fmt.Errorf("发送邮件失败: %w", err)
	}

	// 6. 记录成功历史
	history := s.createEmailHistory(req, models.EmailStatusSuccess, "")
	utils.Infof("邮件发送成功: To=%v, Subject=%s", req.To, req.Subject)

	return history, nil
}

// buildEmailMessage 构建邮件消息
func (s *EmailService) buildEmailMessage(config *models.SMTPConfig, req *SendEmailRequest) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 构建邮件头
	headers := map[string]string{
		"From":         s.formatEmailAddress(config.FromName, config.FromEmail),
		"To":           strings.Join(req.To, ", "),
		"Subject":      req.Subject,
		"Date":         time.Now().Format(time.RFC1123Z),
		"MIME-Version": "1.0",
	}

	// 添加抄送
	if len(req.Cc) > 0 {
		headers["Cc"] = strings.Join(req.Cc, ", ")
	}

	// 如果有附件，使用multipart/mixed
	if len(req.Attachments) > 0 {
		headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary())

		// 写入邮件头
		for k, v := range headers {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
		buf.WriteString("\r\n")

		// 创建HTML正文部分
		htmlHeader := textproto.MIMEHeader{}
		htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
		htmlPart, err := writer.CreatePart(htmlHeader)
		if err != nil {
			return nil, fmt.Errorf("创建HTML部分失败: %w", err)
		}
		htmlPart.Write([]byte(req.Body))

		// 添加附件
		for _, attachment := range req.Attachments {
			if err := s.addAttachment(writer, attachment); err != nil {
				return nil, fmt.Errorf("添加附件失败: %w", err)
			}
		}

		writer.Close()
	} else {
		// 没有附件，直接发送HTML
		headers["Content-Type"] = "text/html; charset=UTF-8"

		// 写入邮件头
		for k, v := range headers {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
		buf.WriteString("\r\n")
		buf.WriteString(req.Body)
	}

	return buf.Bytes(), nil
}

// addAttachment 添加附件
func (s *EmailService) addAttachment(writer *multipart.Writer, attachment Attachment) error {
	utils.Infof("处理附件: Filename=%s, ContentLength=%d", attachment.Filename, len(attachment.Content))
	
	// 解码base64内容
	content, err := base64.StdEncoding.DecodeString(attachment.Content)
	if err != nil {
		utils.Errorf("解码附件内容失败: %v", err)
		return fmt.Errorf("解码附件内容失败: %w", err)
	}

	utils.Infof("附件解码成功: Filename=%s, DecodedSize=%d", attachment.Filename, len(content))

	// 确定内容类型
	contentType := attachment.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 创建附件部分
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", contentType)
	h.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.Filename))
	h.Set("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(h)
	if err != nil {
		return fmt.Errorf("创建附件部分失败: %w", err)
	}

	// 写入base64编码的内容
	encoder := base64.NewEncoder(base64.StdEncoding, part)
	encoder.Write(content)
	encoder.Close()

	return nil
}

// sendEmailViaSMTP 通过SMTP发送邮件
func (s *EmailService) sendEmailViaSMTP(config *models.SMTPConfig, password string, to, cc, bcc []string, message []byte) error {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// 合并所有收件人
	allRecipients := append(append(to, cc...), bcc...)

	// 根据加密类型选择发送方式
	switch config.Encryption {
	case models.EncryptionTLS:
		return s.sendWithTLS(addr, config.Username, password, config.Host, config.FromEmail, allRecipients, message)
	case models.EncryptionStartTLS:
		return s.sendWithStartTLS(addr, config.Username, password, config.Host, config.FromEmail, allRecipients, message)
	default:
		return s.sendPlain(addr, config.Username, password, config.Host, config.FromEmail, allRecipients, message)
	}
}

// sendPlain 普通SMTP发送
func (s *EmailService) sendPlain(addr, username, password, host, from string, to []string, message []byte) error {
	auth := smtp.PlainAuth("", username, password, host)
	return smtp.SendMail(addr, auth, from, to, message)
}

// sendWithTLS 使用TLS发送邮件（端口465 - SMTPS）
func (s *EmailService) sendWithTLS(addr, username, password, host, from string, to []string, message []byte) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	// 使用tls.Dial直接建立加密连接
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS连接失败: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	// 认证
	if username != "" && password != "" {
		auth := smtp.PlainAuth("", username, password, host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("认证失败: %w", err)
		}
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("设置收件人失败 (%s): %w", recipient, err)
		}
	}

	// 发送邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write(message)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	return nil
}

// sendWithStartTLS 使用StartTLS发送（端口587）
func (s *EmailService) sendWithStartTLS(addr, username, password, host, from string, to []string, message []byte) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	// 连接SMTP服务器
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}
	defer client.Close()

	// 检查是否支持STARTTLS
	if ok, _ := client.Extension("STARTTLS"); !ok {
		return errors.New("服务器不支持STARTTLS")
	}

	// 启动STARTTLS
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("启动STARTTLS失败: %w", err)
	}

	// 认证
	if username != "" && password != "" {
		auth := smtp.PlainAuth("", username, password, host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("认证失败: %w", err)
		}
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("设置收件人失败 (%s): %w", recipient, err)
		}
	}

	// 发送邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write(message)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	return nil
}

// createEmailHistory 创建邮件发送历史记录
func (s *EmailService) createEmailHistory(req *SendEmailRequest, status models.EmailStatus, errorMessage string) *models.EmailHistory {
	// 转换附件格式
	attachments := make([]models.Attachment, len(req.Attachments))
	for i, att := range req.Attachments {
		attachments[i] = models.Attachment{
			Filename: att.Filename,
			Path:     "", // base64内容不保存路径
			Size:     int64(len(att.Content)),
		}
	}

	// 使用第一个收件人作为主要收件人
	toEmail := req.To[0]
	if len(req.To) > 1 {
		toEmail = strings.Join(req.To, ", ")
	}

	history := &models.EmailHistory{
		SmtpConfigID: req.SmtpConfigID,
		ToEmail:      toEmail,
		CcEmail:      req.Cc,
		BccEmail:     req.Bcc,
		Subject:      req.Subject,
		Body:         req.Body,
		Attachments:  attachments,
		Status:       status,
		ErrorMessage: errorMessage,
		SentAt:       time.Now(),
	}

	// 保存到数据库
	db := database.GetDB()
	if err := db.Create(history).Error; err != nil {
		utils.Errorf("保存邮件历史失败: %v", err)
	}

	return history
}

// formatEmailAddress 格式化邮箱地址
func (s *EmailService) formatEmailAddress(name, email string) string {
	if name != "" {
		return fmt.Sprintf("%s <%s>", name, email)
	}
	return email
}

// validateEmails 验证邮箱格式
func validateEmails(emails []string) error {
	for _, email := range emails {
		if _, err := mail.ParseAddress(email); err != nil {
			return fmt.Errorf("无效的邮箱地址: %s", email)
		}
	}
	return nil
}
