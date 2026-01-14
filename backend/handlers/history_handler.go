package handlers

import (
	"net/http"
	"strconv"

	"smtp-mail/backend/services"
	"smtp-mail/backend/utils"

	"github.com/gin-gonic/gin"
)

// HistoryHandler 历史处理器
type HistoryHandler struct {
	historyService *services.HistoryService
}

// NewHistoryHandler 创建历史处理器实例
func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{
		historyService: services.NewHistoryService(),
	}
}

// GetAllHistory 获取发送历史（支持分页和状态筛选）
// GET /api/history?page=1&pageSize=10&status=all
func (h *HistoryHandler) GetAllHistory(c *gin.Context) {
	// 解析查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	status := c.DefaultQuery("status", "all")

	// 转换页码
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 转换每页数量
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 验证状态参数
	if status != "all" && status != "success" && status != "failed" {
		errorResponse(c, http.StatusBadRequest, "无效的状态参数，必须是 all/success/failed", nil)
		return
	}

	// 获取历史记录
	result, err := h.historyService.GetAllHistory(page, pageSize, status)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "获取历史记录失败", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", result)
}

// GetHistoryByID 获取单条历史记录
// GET /api/history/:id
func (h *HistoryHandler) GetHistoryByID(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的历史记录ID", err)
		return
	}

	// 获取历史记录
	history, err := h.historyService.GetHistoryByID(uint(id))
	if err != nil {
		errorResponse(c, http.StatusNotFound, "历史记录不存在", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", history)
}

// DeleteHistory 删除历史记录
// DELETE /api/history/:id
func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的历史记录ID", err)
		return
	}

	// 删除历史记录
	if err := h.historyService.DeleteHistory(uint(id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, "删除历史记录失败", err)
		return
	}

	utils.Infof("删除历史记录成功: ID=%d", id)
	successResponse(c, http.StatusOK, "删除成功", nil)
}

// GetStatistics 获取统计信息
// GET /api/history/statistics
func (h *HistoryHandler) GetStatistics(c *gin.Context) {
	// 获取统计信息
	stats, err := h.historyService.GetStatistics()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "获取统计信息失败", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", stats)
}

// RegisterRoutes 注册路由
func (h *HistoryHandler) RegisterRoutes(router *gin.RouterGroup) {
	historyGroup := router.Group("/history")
	{
		historyGroup.GET("", h.GetAllHistory)           // 获取历史记录列表
		historyGroup.GET("/statistics", h.GetStatistics) // 获取统计信息
		historyGroup.GET("/:id", h.GetHistoryByID)       // 获取单条历史记录
		historyGroup.DELETE("/:id", h.DeleteHistory)     // 删除历史记录
	}
}