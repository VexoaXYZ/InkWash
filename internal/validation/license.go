package validation

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	KeyPrefix    = "cfxk_"
	MinKeyLength = 20 // cfxk_ (5) + minimum 15 chars
	MaxKeyLength = 64
)

var keyPattern = regexp.MustCompile(`^cfxk_[A-Za-z0-9_]+$`)

// ValidationError represents a validation error with helpful context
type ValidationError struct {
	Field   string
	Message string
	Hint    string
}

func (e *ValidationError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s: %s\nHint: %s", e.Field, e.Message, e.Hint)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateLicenseKey validates a FiveM license key format
func ValidateLicenseKey(key string) error {
	if key == "" {
		return &ValidationError{
			Field:   "license_key",
			Message: "License key cannot be empty",
			Hint:    "Get your key from https://keymaster.fivem.net",
		}
	}

	if !strings.HasPrefix(key, KeyPrefix) {
		return &ValidationError{
			Field:   "license_key",
			Message: fmt.Sprintf("Must start with '%s'", KeyPrefix),
			Hint:    "Valid format: cfxk_XXXXXXXXXXXX",
		}
	}

	if len(key) < MinKeyLength {
		return &ValidationError{
			Field:   "license_key",
			Message: fmt.Sprintf("Too short (minimum %d chars)", MinKeyLength),
			Hint:    "Check you've copied the complete key",
		}
	}

	if len(key) > MaxKeyLength {
		return &ValidationError{
			Field:   "license_key",
			Message: fmt.Sprintf("Too long (maximum %d chars)", MaxKeyLength),
			Hint:    "Check for extra characters or spaces",
		}
	}

	if !keyPattern.MatchString(key) {
		return &ValidationError{
			Field:   "license_key",
			Message: "Contains invalid characters",
			Hint:    "Only letters, numbers, and underscores allowed after 'cfxk_'",
		}
	}

	keyBody := strings.TrimPrefix(key, KeyPrefix)
	if len(keyBody) < 15 {
		return &ValidationError{
			Field:   "license_key",
			Message: "Key body too short",
			Hint:    "Portion after 'cfxk_' must be at least 15 characters",
		}
	}

	return nil
}

// MaskKey masks a license key for display, showing only prefix and suffix
func MaskKey(key string) string {
	if len(key) < 15 {
		return "****"
	}
	return key[:5] + strings.Repeat("*", len(key)-9) + key[len(key)-4:]
}
