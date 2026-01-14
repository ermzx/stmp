package handlers

import (
	"net/http"
	"strconv"

	"smtp-mail/backend/models"
	"smtp-mail/backend/services"
	"smtp-mail/backend/utils"

	"github.com/gin-gonic/gin"
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	templateService *services.TemplateService
}

// NewTemplateHandler 创建模板处理器实例
func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{
		templateService: services.NewTemplateService(),
	}
}

// GetAllTemplates 获取所有邮件模板
// GET /api/templates
func (h *TemplateHandler) GetAllTemplates(c *gin.Context) {
	templates, err := h.templateService.GetAllTemplates()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "获取模板列表失败", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", templates)
}

// GetTemplateByID 获取单个模板
// GET /api/templates/:id
func (h *TemplateHandler) GetTemplateByID(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的模板ID", err)
		return
	}

	// 获取模板
	template, err := h.templateService.GetTemplateByID(uint(id))
	if err != nil {
		errorResponse(c, http.StatusNotFound, "模板不存在", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", template)
}

// CreateTemplate 创建模板
// POST /api/templates
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var template models.EmailTemplate

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&template); err != nil {
		errorResponse(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 验证必填字段
	if template.Name == "" {
		errorResponse(c, http.StatusBadRequest, "模板名称不能为空", nil)
		return
	}
	if template.Subject == "" {
		errorResponse(c, http.StatusBadRequest, "邮件主题不能为空", nil)
		return
	}
	if template.Body == "" {
		errorResponse(c, http.StatusBadRequest, "邮件正文不能为空", nil)
		return
	}

	// 检查名称是否已存在
	exists, err := h.templateService.NameExists(template.Name, 0)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "检查模板名称失败", err)
		return
	}
	if exists {
		errorResponse(c, http.StatusBadRequest, "模板名称已存在", nil)
		return
	}

	// 创建模板
	if err := h.templateService.CreateTemplate(&template); err != nil {
		// 检查是否是验证错误
		if err == models.ErrTemplateNameRequired || 
		   err == models.ErrTemplateSubjectRequired || 
		   err == models.ErrTemplateBodyRequired {
			errorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "创建模板失败", err)
		return
	}

	utils.Infof("创建模板成功: ID=%d, Name=%s", template.ID, template.Name)
	successResponse(c, http.StatusCreated, "创建成功", template)
}

// UpdateTemplate 更新模板
// PUT /api/templates/:id
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的模板ID", err)
		return
	}

	var template models.EmailTemplate

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&template); err != nil {
		errorResponse(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 验证必填字段
	if template.Name == "" {
		errorResponse(c, http.StatusBadRequest, "模板名称不能为空", nil)
		return
	}
	if template.Subject == "" {
		errorResponse(c, http.StatusBadRequest, "邮件主题不能为空", nil)
		return
	}
	if template.Body == "" {
		errorResponse(c, http.StatusBadRequest, "邮件正文不能为空", nil)
		return
	}

	// 验证更新数据（包括名称唯一性检查）
	if err := h.templateService.ValidateTemplateForUpdate(uint(id), &template); err != nil {
		if err == models.ErrTemplateNameRequired || 
		   err == models.ErrTemplateSubjectRequired || 
		   err == models.ErrTemplateBodyRequired {
			errorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		errorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 更新模板
	if err := h.templateService.UpdateTemplate(uint(id), &template); err != nil {
		errorResponse(c, http.StatusInternalServerError, "更新模板失败", err)
		return
	}

	// 获取更新后的模板
	updatedTemplate, err := h.templateService.GetTemplateByID(uint(id))
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "获取更新后的模板失败", err)
		return
	}

	utils.Infof("更新模板成功: ID=%d", id)
	successResponse(c, http.StatusOK, "更新成功", updatedTemplate)
}

// DeleteTemplate 删除模板
// DELETE /api/templates/:id
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的模板ID", err)
		return
	}

	// 删除模板
	if err := h.templateService.DeleteTemplate(uint(id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, "删除模板失败", err)
		return
	}

	utils.Infof("删除模板成功: ID=%d", id)
	successResponse(c, http.StatusOK, "删除成功", nil)
}

// RegisterRoutes 注册路由
func (h *TemplateHandler) RegisterRoutes(router *gin.RouterGroup) {
	templateGroup := router.Group("/templates")
	{
		templateGroup.GET("", h.GetAllTemplates)       // 获取所有模板
		templateGroup.POST("", h.CreateTemplate)       // 创建模板
		templateGroup.GET("/:id", h.GetTemplateByID)   // 获取单个模板
		templateGroup.PUT("/:id", h.UpdateTemplate)    // 更新模板
		templateGroup.DELETE("/:id", h.DeleteTemplate) // 删除模板
	}
}