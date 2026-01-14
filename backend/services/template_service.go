package services

import (
	"errors"
	"fmt"

	"smtp-mail/backend/database"
	"smtp-mail/backend/models"
	"smtp-mail/backend/utils"
)

// TemplateService 模板服务
type TemplateService struct{}

// NewTemplateService 创建模板服务实例
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// GetAllTemplates 获取所有邮件模板
func (s *TemplateService) GetAllTemplates() ([]models.EmailTemplate, error) {
	db := database.GetDB()
	var templates []models.EmailTemplate

	if err := db.Find(&templates).Error; err != nil {
		utils.Errorf("获取所有模板失败: %v", err)
		return nil, fmt.Errorf("获取所有模板失败: %w", err)
	}

	utils.Infof("获取所有模板成功，共 %d 个", len(templates))
	return templates, nil
}

// GetTemplateByID 获取单个模板
func (s *TemplateService) GetTemplateByID(id uint) (*models.EmailTemplate, error) {
	db := database.GetDB()
	var template models.EmailTemplate

	if err := db.First(&template, id).Error; err != nil {
		utils.Errorf("获取模板失败 (ID: %d): %v", id, err)
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	utils.Infof("获取模板成功: ID=%d, Name=%s", template.ID, template.Name)
	return &template, nil
}

// CreateTemplate 创建模板
func (s *TemplateService) CreateTemplate(template *models.EmailTemplate) error {
	// 验证数据
	if err := s.validateTemplate(template); err != nil {
		return err
	}

	db := database.GetDB()
	if err := db.Create(template).Error; err != nil {
		utils.Errorf("创建模板失败: %v", err)
		return fmt.Errorf("创建模板失败: %w", err)
	}

	utils.Infof("创建模板成功: ID=%d, Name=%s", template.ID, template.Name)
	return nil
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(id uint, template *models.EmailTemplate) error {
	// 验证数据
	if err := s.validateTemplate(template); err != nil {
		return err
	}

	db := database.GetDB()

	// 检查模板是否存在
	var existingTemplate models.EmailTemplate
	if err := db.First(&existingTemplate, id).Error; err != nil {
		utils.Errorf("模板不存在 (ID: %d): %v", id, err)
		return fmt.Errorf("模板不存在: %w", err)
	}

	// 更新模板
	if err := db.Model(&existingTemplate).Updates(template).Error; err != nil {
		utils.Errorf("更新模板失败 (ID: %d): %v", id, err)
		return fmt.Errorf("更新模板失败: %w", err)
	}

	utils.Infof("更新模板成功: ID=%d", id)
	return nil
}

// DeleteTemplate 删除模板
func (s *TemplateService) DeleteTemplate(id uint) error {
	db := database.GetDB()

	// 检查模板是否存在
	var template models.EmailTemplate
	if err := db.First(&template, id).Error; err != nil {
		utils.Errorf("模板不存在 (ID: %d): %v", id, err)
		return fmt.Errorf("模板不存在: %w", err)
	}

	// 删除模板
	if err := db.Delete(&template).Error; err != nil {
		utils.Errorf("删除模板失败 (ID: %d): %v", id, err)
		return fmt.Errorf("删除模板失败: %w", err)
	}

	utils.Infof("删除模板成功: ID=%d, Name=%s", id, template.Name)
	return nil
}

// validateTemplate 验证模板数据
func (s *TemplateService) validateTemplate(template *models.EmailTemplate) error {
	if template.Name == "" {
		return models.ErrTemplateNameRequired
	}
	if template.Subject == "" {
		return models.ErrTemplateSubjectRequired
	}
	if template.Body == "" {
		return models.ErrTemplateBodyRequired
	}
	return nil
}

// GetTemplateByName 根据名称获取模板
func (s *TemplateService) GetTemplateByName(name string) (*models.EmailTemplate, error) {
	db := database.GetDB()
	var template models.EmailTemplate

	if err := db.Where("name = ?", name).First(&template).Error; err != nil {
		utils.Errorf("根据名称获取模板失败 (Name: %s): %v", name, err)
		return nil, fmt.Errorf("根据名称获取模板失败: %w", err)
	}

	utils.Infof("根据名称获取模板成功: ID=%d, Name=%s", template.ID, template.Name)
	return &template, nil
}

// TemplateExists 检查模板是否存在
func (s *TemplateService) TemplateExists(id uint) (bool, error) {
	db := database.GetDB()
	var count int64

	if err := db.Model(&models.EmailTemplate{}).Where("id = ?", id).Count(&count).Error; err != nil {
		utils.Errorf("检查模板是否存在失败 (ID: %d): %v", id, err)
		return false, fmt.Errorf("检查模板是否存在失败: %w", err)
	}

	return count > 0, nil
}

// NameExists 检查模板名称是否已存在
func (s *TemplateService) NameExists(name string, excludeID uint) (bool, error) {
	db := database.GetDB()
	var count int64

	query := db.Model(&models.EmailTemplate{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		utils.Errorf("检查模板名称是否存在失败 (Name: %s): %v", name, err)
		return false, fmt.Errorf("检查模板名称是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ValidateTemplateForUpdate 验证更新时的模板数据（包括名称唯一性检查）
func (s *TemplateService) ValidateTemplateForUpdate(id uint, template *models.EmailTemplate) error {
	// 基本验证
	if err := s.validateTemplate(template); err != nil {
		return err
	}

	// 检查名称是否已被其他模板使用
	exists, err := s.NameExists(template.Name, id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("模板名称已存在")
	}

	return nil
}