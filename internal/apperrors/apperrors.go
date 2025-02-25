package apperrors

import "fmt"

type ConflictError struct {
	ShortURL string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: original URL already exists, short URL: %s", e.ShortURL)
}

func NewConflictError(shortURL string) *ConflictError {
	return &ConflictError{
		ShortURL: shortURL,
	}
}
