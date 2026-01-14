package services

import (
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"smtp-mail/backend/database"
	"smtp-mail/backend/models"
	"smtp-mail/backend/utils"
)

// SMTPService SMTP服务
type SMTPService struct {
	cryptoService *CryptoService
}

// NewSMTPService 创建SMTP服务实例
func NewSMTPService() *SMTPService {
	return &SMTPService{
		cryptoService: NewCryptoService(),
	}
}

// GetAllConfigs 获取所有SMTP配置（不返回密码）
func (s *SMTPService) GetAllConfigs() ([]models.SMTPConfig, error) {
	var configs []models.SMTPConfig
	db := database.GetDB()

	err := db.Find(&configs).Error
	if err != nil {
		utils.Errorf("获取所有SMTP配置失败: %v", err)
		return nil, err
	}

	// 清除密码字段
	for i := range configs {
		configs[i].Password = ""
	}

	utils.Infof("成功获取 %d 个SMTP配置", len(configs))
	return configs, nil
}

// GetConfigByID 获取单个配置（不返回密码）
func (s *SMTPService) GetConfigByID(id uint) (*models.SMTPConfig, error) {
	var config models.SMTPConfig
	db := database.GetDB()

	err := db.First(&config, id).Error
	if err != nil {
		utils.Errorf("获取SMTP配置失败 (ID: %d): %v", id, err)
		return nil, err
	}

	// 清除密码字段
	config.Password = ""

	utils.Infof("成功获取SMTP配置 (ID: %d)", id)
	return &config, nil
}

// GetConfigByIDWithPassword 获取单个配置（包含密码，用于内部使用）
func (s *SMTPService) GetConfigByIDWithPassword(id uint) (*models.SMTPConfig, error) {
	var config models.SMTPConfig
	db := database.GetDB()

	err := db.First(&config, id).Error
	if err != nil {
		utils.Errorf("获取SMTP配置失败 (ID: %d): %v", id, err)
		return nil, err
	}

	utils.Infof("成功获取SMTP配置（含密码）(ID: %d)", id)
	return &config, nil
}

// CreateConfig 创建配置（密码加密）
func (s *SMTPService) CreateConfig(config *models.SMTPConfig) error {
	// 加密密码（使用AES加密，SMTP认证需要可解密的密码）
	if config.Password != "" {
		encryptedPassword, err := s.cryptoService.EncryptPassword(config.Password)
		if err != nil {
			utils.Errorf("加密密码失败: %v", err)
			return fmt.Errorf("加密密码失败: %w", err)
		}
		config.Password = encryptedPassword
	}

	db := database.GetDB()

	err := db.Create(config).Error
	if err != nil {
		utils.Errorf("创建SMTP配置失败: %v", err)
		return err
	}

	// 清除密码字段
	config.Password = ""

	utils.Infof("成功创建SMTP配置 (ID: %d, Name: %s)", config.ID, config.Name)
	return nil
}

// UpdateConfig 更新配置
func (s *SMTPService) UpdateConfig(id uint, config *models.SMTPConfig) error {
	db := database.GetDB()

	// 检查配置是否存在
	var existingConfig models.SMTPConfig
	if err := db.First(&existingConfig, id).Error; err != nil {
		utils.Errorf("SMTP配置不存在 (ID: %d): %v", id, err)
		return err
	}

	// 如果提供了新密码，则加密
	if config.Password != "" {
		encryptedPassword, err := s.cryptoService.EncryptPassword(config.Password)
		if err != nil {
			utils.Errorf("加密密码失败: %v", err)
			return fmt.Errorf("加密密码失败: %w", err)
		}
		config.Password = encryptedPassword
	} else {
		// 如果没有提供新密码，保持原密码
		config.Password = existingConfig.Password
	}

	// 更新配置
	err := db.Model(&existingConfig).Updates(config).Error
	if err != nil {
		utils.Errorf("更新SMTP配置失败 (ID: %d): %v", id, err)
		return err
	}

	// 清除密码字段
	config.Password = ""

	utils.Infof("成功更新SMTP配置 (ID: %d)", id)
	return nil
}

// DeleteConfig 删除配置
func (s *SMTPService) DeleteConfig(id uint) error {
	db := database.GetDB()

	// 检查配置是否存在
	var config models.SMTPConfig
	if err := db.First(&config, id).Error; err != nil {
		utils.Errorf("SMTP配置不存在 (ID: %d): %v", id, err)
		return err
	}

	// 删除配置
	err := db.Delete(&config).Error
	if err != nil {
		utils.Errorf("删除SMTP配置失败 (ID: %d): %v", id, err)
		return err
	}

	utils.Infof("成功删除SMTP配置 (ID: %d)", id)
	return nil
}

// SetDefaultConfig 设置默认配置
func (s *SMTPService) SetDefaultConfig(id uint) error {
	db := database.GetDB()

	// 检查配置是否存在
	var config models.SMTPConfig
	if err := db.First(&config, id).Error; err != nil {
		utils.Errorf("SMTP配置不存在 (ID: %d): %v", id, err)
		return err
	}

	// 开始事务
	tx := db.Begin()

	// 将所有配置的IsDefault设置为false
	if err := tx.Model(&models.SMTPConfig{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		utils.Errorf("重置默认配置失败: %v", err)
		return err
	}

	// 将指定配置设置为默认
	if err := tx.Model(&config).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		utils.Errorf("设置默认配置失败 (ID: %d): %v", id, err)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		utils.Errorf("提交事务失败: %v", err)
		return err
	}

	utils.Infof("成功设置默认SMTP配置 (ID: %d)", id)
	return nil
}

// GetDefaultConfig 获取默认配置
func (s *SMTPService) GetDefaultConfig() (*models.SMTPConfig, error) {
	var config models.SMTPConfig
	db := database.GetDB()

	err := db.Where("is_default = ?", true).First(&config).Error
	if err != nil {
		utils.Errorf("获取默认SMTP配置失败: %v", err)
		return nil, err
	}

	// 清除密码字段
	config.Password = ""

	utils.Infof("成功获取默认SMTP配置 (ID: %d)", config.ID)
	return &config, nil
}

// TestConnection 测试SMTP连接
func (s *SMTPService) TestConnection(config *models.SMTPConfig) error {
	// 解密密码
	password, err := s.cryptoService.DecryptPassword(config.Password)
	if err != nil {
		utils.Errorf("解密密码失败: %v", err)
		return fmt.Errorf("解密密码失败: %w", err)
	}

	// 构建SMTP服务器地址
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// 根据加密类型选择认证方式
	var auth smtp.Auth
	if config.Username != "" && password != "" {
		auth = smtp.PlainAuth("", config.Username, password, config.Host)
	}

	// 测试连接
	var connErr error
	if config.Encryption == models.EncryptionTLS {
		// TLS连接
		connErr = testTLSConnection(addr, auth, config.Host)
	} else if config.Encryption == models.EncryptionStartTLS {
		// StartTLS连接
		connErr = testStartTLSConnection(addr, auth, config.Host)
	} else {
		// 普通连接
		connErr = testPlainConnection(addr, auth)
	}

	if connErr != nil {
		utils.Errorf("SMTP连接测试失败: %v", connErr)
		return fmt.Errorf("连接测试失败: %w", connErr)
	}

	utils.Infof("SMTP连接测试成功 (Host: %s, Port: %d)", config.Host, config.Port)
	return nil
}

// testPlainConnection 测试普通SMTP连接
func testPlainConnection(addr string, auth smtp.Auth) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	return nil
}

// testTLSConnection 测试TLS SMTP连接
func testTLSConnection(addr string, auth smtp.Auth, host string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	// 启动TLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(nil); err != nil {
			return err
		}
	} else {
		return errors.New("服务器不支持STARTTLS")
	}

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	return nil
}

// testStartTLSConnection 测试StartTLS SMTP连接
func testStartTLSConnection(addr string, auth smtp.Auth, host string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	// 启动StartTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(nil); err != nil {
			return err
		}
	} else {
		return errors.New("服务器不支持STARTTLS")
	}

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	return nil
}

// SendTestEmail 发送测试邮件
func (s *SMTPService) SendTestEmail(config *models.SMTPConfig, toEmail string) error {
	// 解密密码
	password, err := s.cryptoService.DecryptPassword(config.Password)
	if err != nil {
		utils.Errorf("解密密码失败: %v", err)
		return fmt.Errorf("解密密码失败: %w", err)
	}

	// 构建SMTP服务器地址
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// 构建邮件内容
	from := config.FromEmail
	if config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail)
	}

	subject := "SMTP配置测试邮件"
	body := fmt.Sprintf("这是一封测试邮件，用于验证SMTP配置是否正确。\n\n配置名称: %s\n发送时间: %s\n\n如果您收到此邮件，说明SMTP配置成功！",
		config.Name, time.Now().Format("2006-01-02 15:04:05"))

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, toEmail, subject, body)

	// 根据加密类型选择发送方式
	var sendErr error
	if config.Encryption == models.EncryptionTLS {
		sendErr = sendTLS(addr, config.Username, password, config.Host, from, []string{toEmail}, []byte(message))
	} else if config.Encryption == models.EncryptionStartTLS {
		sendErr = sendStartTLS(addr, config.Username, password, config.Host, from, []string{toEmail}, []byte(message))
	} else {
		sendErr = smtp.SendMail(addr, smtp.PlainAuth("", config.Username, password, config.Host), from, []string{toEmail}, []byte(message))
	}

	if sendErr != nil {
		utils.Errorf("发送测试邮件失败: %v", sendErr)
		return fmt.Errorf("发送测试邮件失败: %w", sendErr)
	}

	utils.Infof("成功发送测试邮件到 %s", toEmail)
	return nil
}

// sendTLS 使用TLS发送邮件
func sendTLS(addr, username, password, host, from string, to []string, msg []byte) error {
	auth := smtp.PlainAuth("", username, password, host)
	return smtp.SendMail(addr, auth, from, to, msg)
}

// sendStartTLS 使用StartTLS发送邮件
func sendStartTLS(addr, username, password, host, from string, to []string, msg []byte) error {
	auth := smtp.PlainAuth("", username, password, host)
	return smtp.SendMail(addr, auth, from, to, msg)
}
