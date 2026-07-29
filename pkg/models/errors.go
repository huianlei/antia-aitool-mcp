package models

import "fmt"

// Error types
var (
	ErrPluginNotFound    = fmt.Errorf("plugin not found")
	ErrPluginDisabled    = fmt.Errorf("plugin disabled")
	ErrToolNotFound      = fmt.Errorf("tool not found")
	ErrInvalidToolName   = fmt.Errorf("invalid tool name format")
	ErrPluginInitFailed  = fmt.Errorf("plugin initialization failed")
	ErrInvalidConfig     = fmt.Errorf("invalid configuration")
	ErrConnectionFailed  = fmt.Errorf("connection failed")
	ErrAuthFailed        = fmt.Errorf("authentication failed")
	ErrRequestTimeout    = fmt.Errorf("request timeout")
	ErrInvalidParameters = fmt.Errorf("invalid parameters")
)

// AppError represents an application error with context
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
