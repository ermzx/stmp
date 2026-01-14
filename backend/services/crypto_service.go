package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"smtp-mail/backend/config"
	"smtp-mail/backend/utils"
)

// CryptoService 加密服务
type CryptoService struct {
	cost          int
	encryptionKey string
}

// NewCryptoService 创建加密服务实例
func NewCryptoService() *CryptoService {
	cfg := config.GetConfig()
	cost := cfg.Security.BcryptCost
	if cost < 4 {
		cost = 10 // 默认值
	}

	// 使用JWT密钥的SHA256哈希作为AES密钥
	encryptionKey := deriveEncryptionKey(cfg.Security.JWTSecret)

	return &CryptoService{
		cost:          cost,
		encryptionKey: encryptionKey,
	}
}

// deriveEncryptionKey 从密钥派生AES密钥
func deriveEncryptionKey(secret string) string {
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	hash := sha256.Sum256([]byte(secret))
	return string(hash[:])
}

// EncryptPassword 使用AES-GCM加密密码（用于SMTP密码）
func (s *CryptoService) EncryptPassword(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher([]byte(s.encryptionKey))
	if err != nil {
		return "", fmt.Errorf("创建AES密码块失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成nonce失败: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword 解密SMTP密码
func (s *CryptoService) DecryptPassword(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// 尝试base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		// 如果不是有效的base64，可能是明文密码，直接返回
		utils.Infof("密码不是有效的base64格式，使用明文密码")
		return ciphertext, nil
	}

	// 检查解码后的数据是否太短（不足以包含nonce）
	if len(data) < 12 {
		// 数据太短，可能是明文密码
		utils.Infof("解码后的数据太短，使用明文密码")
		return ciphertext, nil
	}

	block, err := aes.NewCipher([]byte(s.encryptionKey))
	if err != nil {
		return "", fmt.Errorf("创建AES密码块失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		// 数据太短，可能是明文密码
		utils.Infof("数据长度小于nonce大小，使用明文密码")
		return ciphertext, nil
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		// 解密失败，可能是明文密码
		utils.Infof("解密失败，使用明文密码: %v", err)
		return ciphertext, nil
	}

	return string(plaintext), nil
}

// HashPassword 加密密码（bcrypt，用于用户密码）
func (s *CryptoService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码
func (s *CryptoService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
