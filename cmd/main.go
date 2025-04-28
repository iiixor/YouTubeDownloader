package main

import (
	"YouTubeDownloader/internal/logger"
	"YouTubeDownloader/internal/service"
	"YouTubeDownloader/internal/usecase"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Инициализируем глобальный логгер
	if err := logger.Init(); err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logger.Lg().Sync()

	logger.Lg().Info("Application is starting")

	// Создаём Gin-роутер
	r := gin.New()

	// Middleware: восстановление после паник
	r.Use(gin.Recovery())

	// Middleware: логируем каждый запрос
	r.Use(func(c *gin.Context) {
		logger.Lg().Info("incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	})

	// Инициализируем слои приложения
	svc := service.NewYoutubeService()
	uc := usecase.NewDownloadUsecase(svc)

	// Регистрируем маршрут скачивания
	r.POST("/download", func(c *gin.Context) {
		stream, filename, appErr := uc.DownloadVideo(c.Request.Context(), c.PostForm("url"))
		if appErr != nil {
			logger.Lg().Error("download error", zap.Int("code", appErr.Code), zap.String("message", appErr.Message))
			c.JSON(appErr.Code, gin.H{"code": appErr.Code, "message": appErr.Message})
			return
		}
		defer stream.Close()

		// Заголовки для скачивания файла
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Type", "video/mp4")
		c.Status(200)

		if _, err := io.Copy(c.Writer, stream); err != nil {
			logger.Lg().Error("streaming error", zap.Error(err))
		}
	})

	addr := ":8080"
	logger.Lg().Info("server listening", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Lg().Fatal("failed to run server", zap.Error(err))
	}
}
