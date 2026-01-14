package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smtp-mail/backend/config"
	"smtp-mail/backend/database"
	"smtp-mail/backend/handlers"
	"smtp-mail/backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg := config.GetConfig()
	log.Printf("配置加载成功: 服务器端口=%d, 模式=%s", cfg.Server.Port, cfg.Server.Mode)

	// 设置Gin运行模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := database.Initialize(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()

	// 创建Gin路由实例
	router := gin.New()

	// 注册中间件
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// 创建处理器实例
	smtpHandler := handlers.NewSMTPHandler()
	emailHandler := handlers.NewEmailHandler()
	templateHandler := handlers.NewTemplateHandler()
	historyHandler := handlers.NewHistoryHandler()

	// 注册健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// 注册API路由
	api := router.Group("/api")
	{
		// SMTP配置管理路由
		smtpHandler.RegisterRoutes(api)

		// 邮件发送路由
		emailHandler.RegisterRoutes(api)

		// 邮件模板管理路由
		templateHandler.RegisterRoutes(api)

		// 发送历史记录路由
		historyHandler.RegisterRoutes(api)
	}

	// 配置静态文件服务
	if _, err := os.Stat("../frontend/dist"); err == nil {
		// 提供静态文件服务
		router.Static("/assets", "../frontend/dist/assets")
		router.StaticFile("/favicon.ico", "../frontend/dist/favicon.ico")

		// 所有非API路由都返回index.html（支持前端路由）
		router.NoRoute(func(c *gin.Context) {
			// 如果是API请求，返回404
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "接口不存在",
				})
				return
			}
			// 否则返回index.html
			c.File("../frontend/dist/index.html")
		})

		log.Println("静态文件服务已启用: ../frontend/dist")
	} else {
		log.Println("前端构建目录不存在，跳过静态文件服务")
	}

	// 创建HTTP服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("服务器启动成功，监听端口: %d", cfg.Server.Port)
		log.Printf("访问地址: http://localhost:%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 设置5秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	log.Println("服务器已退出")
}
