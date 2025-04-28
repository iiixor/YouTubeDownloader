package delivery

import (
	"io"
	"net/http"

	"YouTubeDownloader/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// requestBody описывает формат входного JSON.
type requestBody struct {
	URL string `json:"url" binding:"required"` // ссылка на видео YouTube
}

// Handler хранит зависимости для HTTP-обработок.
type Handler struct {
	uc  usecase.Usecase
	log *zap.Logger
}

// NewHandler создаёт новый HTTP-обработчик.
func NewHandler(uc usecase.Usecase, log *zap.Logger) *Handler {
	return &Handler{uc: uc, log: log}
}

// DownloadVideoHandler обрабатывает POST /download.
// Комментарии и документация - на русском, логи - на английском.
func (h *Handler) DownloadVideoHandler(c *gin.Context) {
	var req requestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid request body"})
		return
	}

	stream, filename, err := h.uc.DownloadVideo(c.Request.Context(), req.URL)
	if err != nil {
		h.log.Error("download error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	defer stream.Close()

	// Заголовки для корректного скачивания на клиенте
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Type", "video/mp4")
	c.Status(http.StatusOK)

	// Передаём поток напрямую в ответ
	if _, err := io.Copy(c.Writer, stream); err != nil {
		h.log.Error("streaming error", zap.Error(err))
	}
}
