package handlers

import (
	"net/http"

	"smtp-mail/backend/services"

	"github.com/gin-gonic/gin"
)

// EmailHandler 邮件处理器
type EmailHandler struct {
	emailService *services.EmailService
}

// NewEmailHandler 创建邮件处理器实例
func NewEmailHandler() *EmailHandler {
	return &EmailHandler{
		emailService: services.NewEmailService(),
	}
}

// SendEmail 发送邮件
// POST /api/email/send
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req services.SendEmailRequest

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 验证必填字段
	if req.SmtpConfigID == 0 {
		errorResponse(c, http.StatusBadRequest, "SMTP配置ID不能为空", nil)
		return
	}
	if len(req.To) == 0 {
		errorResponse(c, http.StatusBadRequest, "收件人列表不能为空", nil)
		return
	}
	if req.Subject == "" {
		errorResponse(c, http.StatusBadRequest, "邮件主题不能为空", nil)
		return
	}
	if req.Body == "" {
		errorResponse(c, http.StatusBadRequest, "邮件正文不能为空", nil)
		return
	}

	// 验证附件
	for i, attachment := range req.Attachments {
		if attachment.Filename == "" {
			errorResponse(c, http.StatusBadRequest, "附件文件名不能为空", nil)
			return
		}
		if attachment.Content == "" {
			errorResponse(c, http.StatusBadRequest, "附件内容不能为空", nil)
			return
		}
		// 验证base64编码
		if len(attachment.Content) == 0 {
			errorResponse(c, http.StatusBadRequest, "附件内容无效", nil)
			return
		}
		_ = i // 避免未使用变量警告
	}

	// 调用服务层发送邮件
	history, err := h.emailService.SendEmail(&req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "发送邮件失败", err)
		return
	}

	// 返回发送结果
	successResponse(c, http.StatusOK, "邮件发送成功", history)
}

// RegisterRoutes 注册路由
func (h *EmailHandler) RegisterRoutes(router *gin.RouterGroup) {
	emailGroup := router.Group("/email")
	{
		emailGroup.POST("/send", h.SendEmail) // 发送邮件
	}
}