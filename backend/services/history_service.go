package services

import (
	"fmt"

	"smtp-mail/backend/database"
	"smtp-mail/backend/models"
	"smtp-mail/backend/utils"
)

// HistoryService 历史服务
type HistoryService struct{}

// NewHistoryService 创建历史服务实例
func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

// HistoryListResponse 历史列表响应
type HistoryListResponse struct {
	List     []models.EmailHistory `json:"list"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	Total  int64 `json:"total"`
	Success int64 `json:"success"`
	Failed int64 `json:"failed"`
}

// GetAllHistory 获取发送历史（支持分页和状态筛选）
func (s *HistoryService) GetAllHistory(page, pageSize int, status string) (*HistoryListResponse, error) {
	db := database.GetDB()
	
	// 设置默认值
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 构建查询
	query := db.Model(&models.EmailHistory{})
	
	// 按状态筛选
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Errorf("获取历史记录总数失败: %v", err)
		return nil, fmt.Errorf("获取历史记录总数失败: %w", err)
	}

	// 分页查询
	var histories []models.EmailHistory
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		utils.Errorf("获取历史记录列表失败: %v", err)
		return nil, fmt.Errorf("获取历史记录列表失败: %w", err)
	}

	utils.Infof("获取历史记录成功，共 %d 条，当前页 %d 条", total, len(histories))
	return &HistoryListResponse{
		List:     histories,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetHistoryByID 获取单条历史记录
func (s *HistoryService) GetHistoryByID(id uint) (*models.EmailHistory, error) {
	db := database.GetDB()
	var history models.EmailHistory

	if err := db.First(&history, id).Error; err != nil {
		utils.Errorf("获取历史记录失败 (ID: %d): %v", id, err)
		return nil, fmt.Errorf("获取历史记录失败: %w", err)
	}

	utils.Infof("获取历史记录成功: ID=%d, ToEmail=%s", history.ID, history.ToEmail)
	return &history, nil
}

// DeleteHistory 删除历史记录
func (s *HistoryService) DeleteHistory(id uint) error {
	db := database.GetDB()

	// 检查历史记录是否存在
	var history models.EmailHistory
	if err := db.First(&history, id).Error; err != nil {
		utils.Errorf("历史记录不存在 (ID: %d): %v", id, err)
		return fmt.Errorf("历史记录不存在: %w", err)
	}

	// 删除历史记录
	if err := db.Delete(&history).Error; err != nil {
		utils.Errorf("删除历史记录失败 (ID: %d): %v", id, err)
		return fmt.Errorf("删除历史记录失败: %w", err)
	}

	utils.Infof("删除历史记录成功: ID=%d", id)
	return nil
}

// GetStatistics 获取统计信息（总数、成功数、失败数）
func (s *HistoryService) GetStatistics() (*StatisticsResponse, error) {
	db := database.GetDB()

	// 获取总数
	var total int64
	if err := db.Model(&models.EmailHistory{}).Count(&total).Error; err != nil {
		utils.Errorf("获取历史记录总数失败: %v", err)
		return nil, fmt.Errorf("获取历史记录总数失败: %w", err)
	}

	// 获取成功数
	var success int64
	if err := db.Model(&models.EmailHistory{}).Where("status = ?", models.EmailStatusSuccess).Count(&success).Error; err != nil {
		utils.Errorf("获取成功记录数失败: %v", err)
		return nil, fmt.Errorf("获取成功记录数失败: %w", err)
	}

	// 获取失败数
	var failed int64
	if err := db.Model(&models.EmailHistory{}).Where("status = ?", models.EmailStatusFailed).Count(&failed).Error; err != nil {
		utils.Errorf("获取失败记录数失败: %v", err)
		return nil, fmt.Errorf("获取失败记录数失败: %w", err)
	}

	utils.Infof("获取统计信息成功: Total=%d, Success=%d, Failed=%d", total, success, failed)
	return &StatisticsResponse{
		Total:  total,
		Success: success,
		Failed: failed,
	}, nil
}