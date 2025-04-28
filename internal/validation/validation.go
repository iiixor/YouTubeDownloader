package validation

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrEmptyURL   = errors.New("url is empty")
	ErrInvalidURL = errors.New("invalid URL format")
	ErrNotYouTube = errors.New("URL is not a YouTube link")
)

// ValidateVideoURL проверяет, что rawURL не пустая строка,
// имеет корректный формат и принадлежит YouTube.
func ValidateVideoURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return ErrEmptyURL
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}
	host := parsed.Hostname()
	if !strings.Contains(host, "youtube.com") && !strings.Contains(host, "youtu.be") {
		return ErrNotYouTube
	}
	return nil
}
