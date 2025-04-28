package service

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/kkdai/youtube/v2"
)

// Service - интерфейс для скачивания видео.
type Service interface {
	Download(ctx context.Context, videoURL string) (io.ReadCloser, string, error)
}

type youtubeService struct {
	client *youtube.Client
}

// NewYoutubeService создаёт реализацию Service на основе youtube.Client.
func NewYoutubeService() Service {
	return &youtubeService{client: &youtube.Client{}}
}

// Download выбирает mp4-формат 720p и отдаёт поток для чтения.
func (s *youtubeService) Download(ctx context.Context, videoURL string) (io.ReadCloser, string, error) {
	video, err := s.client.GetVideoContext(ctx, videoURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get video info: %w", err)
	}

	// Фильтруем только mp4-потоки
	var formats []youtube.Format
	for _, f := range video.Formats {
		if strings.HasPrefix(f.MimeType, "video/mp4") {
			formats = append(formats, f)
		}
	}
	if len(formats) == 0 {
		return nil, "", fmt.Errorf("no mp4 formats found")
	}

	// Сортируем по высоте (Height), чтобы выбрать ближайший к 720p
	sort.Slice(formats, func(i, j int) bool {
		return formats[i].Height < formats[j].Height
	})

	// Ищем формат ровно 720p, иначе ближайший ниже 720
	selected := formats[0]
	for _, f := range formats {
		if f.Height == 720 {
			selected = f
			break
		}
		if f.Height > selected.Height && f.Height < 720 {
			selected = f
		}
	}

	stream, _, err := s.client.GetStreamContext(ctx, video, &selected)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get stream: %w", err)
	}

	filename := fmt.Sprintf("%s.mp4", video.ID)
	return stream, filename, nil
}
