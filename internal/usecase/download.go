package usecase

import (
	"context"
	"io"

	"YouTubeDownloader/internal/errors"
	"YouTubeDownloader/internal/service"
	"YouTubeDownloader/internal/validation"
)

// Usecase - интерфейс бизнес-логики.
type Usecase interface {
	// DownloadVideo валидирует URL и возвращает поток и имя файла или AppError.
	DownloadVideo(ctx context.Context, rawURL string) (io.ReadCloser, string, *errors.AppError)
}

type downloadUsecase struct {
	svc service.Service
}

// NewDownloadUsecase создаёт экземпляр Usecase.
func NewDownloadUsecase(svc service.Service) Usecase {
	return &downloadUsecase{svc: svc}
}

func (u *downloadUsecase) DownloadVideo(ctx context.Context, rawURL string) (io.ReadCloser, string, *errors.AppError) {
	// Валидируем входную ссылку
	if err := validation.ValidateVideoURL(rawURL); err != nil {
		// некорректный URL - возвращаем AppError с кодом 400
		return nil, "", errors.NewAppError(400, err.Error())
	}
	// Запускаем скачивание в сервисе
	stream, filename, err := u.svc.Download(ctx, rawURL)
	if err != nil {
		// ошибка на стороне сервиса - тоже 400
		return nil, "", errors.NewAppError(400, err.Error())
	}
	return stream, filename, nil
}
