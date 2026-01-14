package handlers

import (
	"net/http"
	"strconv"

	"smtp-mail/backend/models"
	"smtp-mail/backend/services"
	"smtp-mail/backend/utils"

	"github.com/gin-gonic/gin"
)

// SMTPHandler SMTP处理器
type SMTPHandler struct {
	smtpService *services.SMTPService
}

// NewSMTPHandler 创建SMTP处理器实例
func NewSMTPHandler() *SMTPHandler {
	return &SMTPHandler{
		smtpService: services.NewSMTPService(),
	}
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// successResponse 成功响应
func successResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    200,  // 统一返回 200 表示业务成功
		Message: message,
		Data:    data,
	})
}

// errorResponse 错误响应
func errorResponse(c *gin.Context, code int, message string, err error) {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	c.JSON(code, response)
}

// GetAllConfigs 获取所有SMTP配置
// GET /api/smtp/configs
func (h *SMTPHandler) GetAllConfigs(c *gin.Context) {
	configs, err := h.smtpService.GetAllConfigs()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "获取SMTP配置失败", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", configs)
}

// GetConfigByID 获取单个SMTP配置
// GET /api/smtp/configs/:id
func (h *SMTPHandler) GetConfigByID(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的配置ID", err)
		return
	}

	// 获取配置
	config, err := h.smtpService.GetConfigByID(uint(id))
	if err != nil {
		errorResponse(c, http.StatusNotFound, "SMTP配置不存在", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", config)
}

// CreateConfig 创建SMTP配置
// POST /api/smtp/configs
func (h *SMTPHandler) CreateConfig(c *gin.Context) {
	var config models.SMTPConfig

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 验证必填字段
	if config.Name == "" {
		errorResponse(c, http.StatusBadRequest, "配置名称不能为空", nil)
		return
	}
	if config.Host == "" {
		errorResponse(c, http.StatusBadRequest, "SMTP服务器地址不能为空", nil)
		return
	}
	if config.Port <= 0 || config.Port > 65535 {
		errorResponse(c, http.StatusBadRequest, "SMTP端口无效", nil)
		return
	}
	if config.FromEmail == "" {
		errorResponse(c, http.StatusBadRequest, "发件人邮箱不能为空", nil)
		return
	}

	// 验证加密类型
	if config.Encryption != models.EncryptionNone && 
	   config.Encryption != models.EncryptionTLS && 
	   config.Encryption != models.EncryptionStartTLS {
		errorResponse(c, http.StatusBadRequest, "无效的加密类型", nil)
		return
	}

	// 创建配置
	if err := h.smtpService.CreateConfig(&config); err != nil {
		utils.Errorf("创建SMTP配置失败: %v", err)
		errorResponse(c, http.StatusInternalServerError, "创建SMTP配置失败", err)
		return
	}

	utils.Infof("创建SMTP配置成功: ID=%d, Name=%s, IsDefault=%v", config.ID, config.Name, config.IsDefault)
	successResponse(c, http.StatusCreated, "创建成功", config)
}

// UpdateConfig 更新SMTP配置
// PUT /api/smtp/configs/:id
func (h *SMTPHandler) UpdateConfig(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的配置ID", err)
		return
	}

	var config models.SMTPConfig

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 验证必填字段
	if config.Name == "" {
		errorResponse(c, http.StatusBadRequest, "配置名称不能为空", nil)
		return
	}
	if config.Host == "" {
		errorResponse(c, http.StatusBadRequest, "SMTP服务器地址不能为空", nil)
		return
	}
	if config.Port <= 0 || config.Port > 65535 {
		errorResponse(c, http.StatusBadRequest, "SMTP端口无效", nil)
		return
	}
	if config.FromEmail == "" {
		errorResponse(c, http.StatusBadRequest, "发件人邮箱不能为空", nil)
		return
	}

	// 验证加密类型
	if config.Encryption != models.EncryptionNone && 
	   config.Encryption != models.EncryptionTLS && 
	   config.Encryption != models.EncryptionStartTLS {
		errorResponse(c, http.StatusBadRequest, "无效的加密类型", nil)
		return
	}

	// 更新配置
	if err := h.smtpService.UpdateConfig(uint(id), &config); err != nil {
		errorResponse(c, http.StatusInternalServerError, "更新SMTP配置失败", err)
		return
	}

	// 获取更新后的配置
	updatedConfig, err := h.smtpService.GetConfigByID(uint(id))
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "获取更新后的配置失败", err)
		return
	}

	utils.Infof("更新SMTP配置成功: ID=%d", id)
	successResponse(c, http.StatusOK, "更新成功", updatedConfig)
}

// DeleteConfig 删除SMTP配置
// DELETE /api/smtp/configs/:id
func (h *SMTPHandler) DeleteConfig(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的配置ID", err)
		return
	}

	// 删除配置
	if err := h.smtpService.DeleteConfig(uint(id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, "删除SMTP配置失败", err)
		return
	}

	utils.Infof("删除SMTP配置成功: ID=%d", id)
	successResponse(c, http.StatusOK, "删除成功", nil)
}

// TestConnection 测试SMTP连接
// POST /api/smtp/configs/:id/test
func (h *SMTPHandler) TestConnection(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的配置ID", err)
		return
	}

	// 获取配置
	config, err := h.smtpService.GetConfigByID(uint(id))
	if err != nil {
		errorResponse(c, http.StatusNotFound, "SMTP配置不存在", err)
		return
	}

	// 从请求体获取密码（如果需要重新测试）
	var requestData struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&requestData); err == nil && requestData.Password != "" {
		config.Password = requestData.Password
	}

	// 测试连接
	if err := h.smtpService.TestConnection(config); err != nil {
		errorResponse(c, http.StatusBadRequest, "SMTP连接测试失败", err)
		return
	}

	utils.Infof("SMTP连接测试成功: ID=%d", id)
	successResponse(c, http.StatusOK, "连接测试成功", nil)
}

// SetDefaultConfig 设置默认SMTP配置
// POST /api/smtp/configs/:id/default
func (h *SMTPHandler) SetDefaultConfig(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的配置ID", err)
		return
	}

	// 设置默认配置
	if err := h.smtpService.SetDefaultConfig(uint(id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, "设置默认配置失败", err)
		return
	}

	utils.Infof("设置默认SMTP配置成功: ID=%d", id)
	successResponse(c, http.StatusOK, "设置成功", nil)
}

// GetDefaultConfig 获取默认SMTP配置
// GET /api/smtp/configs/default
func (h *SMTPHandler) GetDefaultConfig(c *gin.Context) {
	config, err := h.smtpService.GetDefaultConfig()
	if err != nil {
		errorResponse(c, http.StatusNotFound, "未找到默认SMTP配置", err)
		return
	}

	successResponse(c, http.StatusOK, "获取成功", config)
}

// SendTestEmail 发送测试邮件
// POST /api/smtp/configs/:id/send-test
func (h *SMTPHandler) SendTestEmail(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的配置ID", err)
		return
	}

	// 获取配置
	config, err := h.smtpService.GetConfigByID(uint(id))
	if err != nil {
		errorResponse(c, http.StatusNotFound, "SMTP配置不存在", err)
		return
	}

	// 从请求体获取收件人邮箱和密码
	var requestData struct {
		ToEmail string `json:"to_email" binding:"required"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		errorResponse(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 如果提供了密码，使用提供的密码
	if requestData.Password != "" {
		config.Password = requestData.Password
	}

	// 发送测试邮件
	if err := h.smtpService.SendTestEmail(config, requestData.ToEmail); err != nil {
		errorResponse(c, http.StatusBadRequest, "发送测试邮件失败", err)
		return
	}

	utils.Infof("发送测试邮件成功: ID=%d, To=%s", id, requestData.ToEmail)
	successResponse(c, http.StatusOK, "测试邮件发送成功", nil)
}

// RegisterRoutes 注册路由
func (h *SMTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	smtpGroup := router.Group("/smtp")
	{
		configs := smtpGroup.Group("/configs")
		{
			configs.GET("", h.GetAllConfigs)           // 获取所有配置
			configs.GET("/default", h.GetDefaultConfig) // 获取默认配置
			configs.POST("", h.CreateConfig)           // 创建配置
			configs.GET("/:id", h.GetConfigByID)       // 获取单个配置
			configs.PUT("/:id", h.UpdateConfig)        // 更新配置
			configs.DELETE("/:id", h.DeleteConfig)     // 删除配置
			configs.POST("/:id/test", h.TestConnection) // 测试连接
			configs.POST("/:id/default", h.SetDefaultConfig) // 设置为默认
			configs.POST("/:id/send-test", h.SendTestEmail) // 发送测试邮件
		}
	}
}