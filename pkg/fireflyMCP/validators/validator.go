package validators

import (
	"fmt"
	"time"
)

// Validator is the main interface for all validators
type Validator interface {
	Validate() error
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

// DateValidator validates date strings
type DateValidator struct {
	Value    string
	Field    string
	Required bool
	Format   string
}

// Validate checks if the date is valid
func (v DateValidator) Validate() error {
	if v.Value == "" {
		if v.Required {
			return ValidationError{
				Field:   v.Field,
				Message: "date is required",
			}
		}
		return nil
	}

	format := v.Format
	if format == "" {
		format = "2006-01-02"
	}

	_, err := time.Parse(format, v.Value)
	if err != nil {
		return ValidationError{
			Field:   v.Field,
			Message: fmt.Sprintf("invalid date format, expected %s", format),
		}
	}

	return nil
}

// PaginationValidator validates pagination parameters
type PaginationValidator struct {
	Page  int
	Limit int
}

// Validate checks if pagination parameters are valid
func (v PaginationValidator) Validate() error {
	if v.Page < 0 {
		return ValidationError{
			Field:   "page",
			Message: "page must be >= 0",
		}
	}

	if v.Limit < 0 {
		return ValidationError{
			Field:   "limit",
			Message: "limit must be >= 0",
		}
	}

	if v.Limit > 1000 {
		return ValidationError{
			Field:   "limit",
			Message: "limit must be <= 1000",
		}
	}

	return nil
}

// DateRangeValidator validates date range
type DateRangeValidator struct {
	StartDate string
	EndDate   string
	Required  bool
}

// Validate checks if the date range is valid
func (v DateRangeValidator) Validate() error {
	// Validate start date
	if v.StartDate != "" || v.Required {
		startValidator := DateValidator{
			Value:    v.StartDate,
			Field:    "start_date",
			Required: v.Required,
		}
		if err := startValidator.Validate(); err != nil {
			return err
		}
	}

	// Validate end date
	if v.EndDate != "" || v.Required {
		endValidator := DateValidator{
			Value:    v.EndDate,
			Field:    "end_date",
			Required: v.Required,
		}
		if err := endValidator.Validate(); err != nil {
			return err
		}
	}

	// Check if end date is after start date
	if v.StartDate != "" && v.EndDate != "" {
		start, _ := time.Parse("2006-01-02", v.StartDate)
		end, _ := time.Parse("2006-01-02", v.EndDate)
		
		if end.Before(start) {
			return ValidationError{
				Field:   "date_range",
				Message: "end date must be after or equal to start date",
			}
		}
	}

	return nil
}

// RequiredFieldValidator validates required fields
type RequiredFieldValidator struct {
	Value interface{}
	Field string
}

// Validate checks if the required field is present
func (v RequiredFieldValidator) Validate() error {
	if v.Value == nil {
		return ValidationError{
			Field:   v.Field,
			Message: "field is required",
		}
	}

	switch val := v.Value.(type) {
	case string:
		if val == "" {
			return ValidationError{
				Field:   v.Field,
				Message: "field cannot be empty",
			}
		}
	case int, int32, int64, float32, float64:
		// Numbers are valid if not nil
		return nil
	default:
		// For other types, non-nil is valid
		return nil
	}

	return nil
}

// CompositeValidator allows multiple validators to be run together
type CompositeValidator struct {
	Validators []Validator
}

// Validate runs all validators and returns the first error
func (v CompositeValidator) Validate() error {
	for _, validator := range v.Validators {
		if err := validator.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateAll runs all validators and returns all errors
func (v CompositeValidator) ValidateAll() []error {
	var errors []error
	for _, validator := range v.Validators {
		if err := validator.Validate(); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// TransactionTypeValidator validates transaction types
type TransactionTypeValidator struct {
	Type string
}

// Validate checks if the transaction type is valid
func (v TransactionTypeValidator) Validate() error {
	validTypes := map[string]bool{
		"":              true, // Empty is valid (no filter)
		"all":           true,
		"withdrawal":    true,
		"withdrawals":   true,
		"expense":       true,
		"deposit":       true,
		"deposits":      true,
		"income":        true,
		"transfer":      true,
		"transfers":     true,
		"opening_balance": true,
		"reconciliation": true,
	}

	if !validTypes[v.Type] {
		return ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("invalid transaction type: %s", v.Type),
		}
	}

	return nil
}

// AmountValidator validates amount values
type AmountValidator struct {
	Amount   string
	Field    string
	Required bool
}

// Validate checks if the amount is valid
func (v AmountValidator) Validate() error {
	if v.Amount == "" {
		if v.Required {
			return ValidationError{
				Field:   v.Field,
				Message: "amount is required",
			}
		}
		return nil
	}

	// Basic validation - could be enhanced with proper decimal parsing
	// For now, we just check it's not empty if required
	return nil
}

// CombineErrors combines multiple validation errors into a single error
func CombineErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	if len(errors) == 1 {
		return errors[0]
	}

	var messages string
	for i, err := range errors {
		if i > 0 {
			messages += "; "
		}
		messages += err.Error()
	}

	return fmt.Errorf("%s", messages)
}