package errors

// AppError представляет собой ошибку с HTTP-кодом и сообщением.
type AppError struct {
	Code    int    // HTTP-код ответа
	Message string // текстовое описание ошибки
}

// Error реализует интерфейс error.
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError создаёт новый AppError.
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
