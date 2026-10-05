// Package apperror berisi error aplikasi yang membawa status HTTP dan pesan aman
// untuk client. Service mengembalikan *AppError; ErrorHandler di middleware
// yang mengubahnya menjadi response JSON.
package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int                 // status HTTP
	Message string              // pesan aman untuk client
	Errors  map[string][]string // detail validasi per field (opsional)
	Cause   error               // penyebab asli, hanya untuk log, tidak dikirim ke client
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Cause }

func New(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

// Validation membuat error 422 dengan rincian per field.
func Validation(errs map[string][]string) *AppError {
	return &AppError{Status: http.StatusUnprocessableEntity, Message: "Validasi gagal", Errors: errs}
}

// FieldError membuat error validasi 422 untuk satu field.
func FieldError(field, message string) *AppError {
	return Validation(map[string][]string{field: {message}})
}

func Unauthorized(message string) *AppError { return New(http.StatusUnauthorized, message) }
func Forbidden(message string) *AppError    { return New(http.StatusForbidden, message) }
func NotFound(message string) *AppError     { return New(http.StatusNotFound, message) }
func Conflict(message string) *AppError     { return New(http.StatusConflict, message) }
func Unprocessable(message string) *AppError {
	return New(http.StatusUnprocessableEntity, message)
}
func TooManyRequests(message string) *AppError {
	return New(http.StatusTooManyRequests, message)
}

// Internal membungkus error tak terduga. Pesan ke client selalu generik.
func Internal(cause error) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Message: "Terjadi kesalahan pada server",
		Cause:   cause,
	}
}