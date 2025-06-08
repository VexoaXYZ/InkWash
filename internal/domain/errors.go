package domain

import "fmt"

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "not_found"
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeNetwork       ErrorType = "network"
	ErrorTypeFilesystem    ErrorType = "filesystem"
	ErrorTypeConfiguration ErrorType = "configuration"
	ErrorTypePermission    ErrorType = "permission"
	ErrorTypeConflict      ErrorType = "conflict"
	ErrorTypeInternal      ErrorType = "internal"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Type    ErrorType
	Message string
	Details map[string]interface{}
	Cause   error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewError creates a new domain error
func NewError(errType ErrorType, message string) *DomainError {
	return &DomainError{
		Type:    errType,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithCause adds a cause to the error
func (e *DomainError) WithCause(cause error) *DomainError {
	e.Cause = cause
	return e
}

// WithDetail adds a detail to the error
func (e *DomainError) WithDetail(key string, value interface{}) *DomainError {
	e.Details[key] = value
	return e
}

// Common error constructors

// ErrServerNotFound creates a server not found error
func ErrServerNotFound(serverID string) *DomainError {
	return NewError(ErrorTypeNotFound, "server not found").
		WithDetail("server_id", serverID)
}

// ErrResourceNotFound creates a resource not found error
func ErrResourceNotFound(resourceName string) *DomainError {
	return NewError(ErrorTypeNotFound, "resource not found").
		WithDetail("resource_name", resourceName)
}

// ErrInvalidServerConfig creates an invalid server configuration error
func ErrInvalidServerConfig(reason string) *DomainError {
	return NewError(ErrorTypeValidation, "invalid server configuration").
		WithDetail("reason", reason)
}

// ErrServerAlreadyExists creates a server already exists error
func ErrServerAlreadyExists(serverName string) *DomainError {
	return NewError(ErrorTypeConflict, "server already exists").
		WithDetail("server_name", serverName)
}

// ErrDownloadFailed creates a download failed error
func ErrDownloadFailed(url string, cause error) *DomainError {
	return NewError(ErrorTypeNetwork, "download failed").
		WithDetail("url", url).
		WithCause(cause)
}

// ErrPermissionDenied creates a permission denied error
func ErrPermissionDenied(action string) *DomainError {
	return NewError(ErrorTypePermission, "permission denied").
		WithDetail("action", action)
}

// ErrFilesystemOperation creates a filesystem operation error
func ErrFilesystemOperation(operation string, path string, cause error) *DomainError {
	return NewError(ErrorTypeFilesystem, fmt.Sprintf("filesystem operation failed: %s", operation)).
		WithDetail("path", path).
		WithCause(cause)
}