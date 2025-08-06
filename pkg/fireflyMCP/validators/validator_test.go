package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator DateValidator
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid date",
			validator: DateValidator{
				Value: "2024-01-15",
				Field: "test_date",
			},
			wantErr: false,
		},
		{
			name: "empty date when not required",
			validator: DateValidator{
				Value:    "",
				Field:    "test_date",
				Required: false,
			},
			wantErr: false,
		},
		{
			name: "empty date when required",
			validator: DateValidator{
				Value:    "",
				Field:    "test_date",
				Required: true,
			},
			wantErr: true,
			errMsg:  "validation error in field 'test_date': date is required",
		},
		{
			name: "invalid date format",
			validator: DateValidator{
				Value: "15-01-2024",
				Field: "test_date",
			},
			wantErr: true,
			errMsg:  "validation error in field 'test_date': invalid date format, expected 2006-01-02",
		},
		{
			name: "custom format valid",
			validator: DateValidator{
				Value:  "15:30:45",
				Field:  "test_time",
				Format: "15:04:05",
			},
			wantErr: false,
		},
		{
			name: "custom format invalid",
			validator: DateValidator{
				Value:  "25:70:80",
				Field:  "test_time",
				Format: "15:04:05",
			},
			wantErr: true,
			errMsg:  "validation error in field 'test_time': invalid date format, expected 15:04:05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPaginationValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator PaginationValidator
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid pagination",
			validator: PaginationValidator{
				Page:  1,
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "zero values valid",
			validator: PaginationValidator{
				Page:  0,
				Limit: 0,
			},
			wantErr: false,
		},
		{
			name: "negative page",
			validator: PaginationValidator{
				Page:  -1,
				Limit: 10,
			},
			wantErr: true,
			errMsg:  "validation error in field 'page': page must be >= 0",
		},
		{
			name: "negative limit",
			validator: PaginationValidator{
				Page:  1,
				Limit: -10,
			},
			wantErr: true,
			errMsg:  "validation error in field 'limit': limit must be >= 0",
		},
		{
			name: "limit exceeds maximum",
			validator: PaginationValidator{
				Page:  1,
				Limit: 1001,
			},
			wantErr: true,
			errMsg:  "validation error in field 'limit': limit must be <= 1000",
		},
		{
			name: "limit at maximum",
			validator: PaginationValidator{
				Page:  1,
				Limit: 1000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDateRangeValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator DateRangeValidator
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid date range",
			validator: DateRangeValidator{
				StartDate: "2024-01-01",
				EndDate:   "2024-01-31",
			},
			wantErr: false,
		},
		{
			name: "empty dates when not required",
			validator: DateRangeValidator{
				StartDate: "",
				EndDate:   "",
				Required:  false,
			},
			wantErr: false,
		},
		{
			name: "empty dates when required",
			validator: DateRangeValidator{
				StartDate: "",
				EndDate:   "",
				Required:  true,
			},
			wantErr: true,
			errMsg:  "validation error in field 'start_date': date is required",
		},
		{
			name: "end date before start date",
			validator: DateRangeValidator{
				StartDate: "2024-01-31",
				EndDate:   "2024-01-01",
			},
			wantErr: true,
			errMsg:  "validation error in field 'date_range': end date must be after or equal to start date",
		},
		{
			name: "same start and end date",
			validator: DateRangeValidator{
				StartDate: "2024-01-15",
				EndDate:   "2024-01-15",
			},
			wantErr: false,
		},
		{
			name: "invalid start date format",
			validator: DateRangeValidator{
				StartDate: "31-01-2024",
				EndDate:   "2024-01-31",
			},
			wantErr: true,
			errMsg:  "validation error in field 'start_date': invalid date format, expected 2006-01-02",
		},
		{
			name: "invalid end date format",
			validator: DateRangeValidator{
				StartDate: "2024-01-01",
				EndDate:   "31-01-2024",
			},
			wantErr: true,
			errMsg:  "validation error in field 'end_date': invalid date format, expected 2006-01-02",
		},
		{
			name: "only start date",
			validator: DateRangeValidator{
				StartDate: "2024-01-01",
				EndDate:   "",
			},
			wantErr: false,
		},
		{
			name: "only end date",
			validator: DateRangeValidator{
				StartDate: "",
				EndDate:   "2024-01-31",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequiredFieldValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator RequiredFieldValidator
		wantErr   bool
		errMsg    string
	}{
		{
			name: "non-empty string",
			validator: RequiredFieldValidator{
				Value: "test",
				Field: "test_field",
			},
			wantErr: false,
		},
		{
			name: "empty string",
			validator: RequiredFieldValidator{
				Value: "",
				Field: "test_field",
			},
			wantErr: true,
			errMsg:  "validation error in field 'test_field': field cannot be empty",
		},
		{
			name: "nil value",
			validator: RequiredFieldValidator{
				Value: nil,
				Field: "test_field",
			},
			wantErr: true,
			errMsg:  "validation error in field 'test_field': field is required",
		},
		{
			name: "integer value",
			validator: RequiredFieldValidator{
				Value: 42,
				Field: "test_field",
			},
			wantErr: false,
		},
		{
			name: "zero integer",
			validator: RequiredFieldValidator{
				Value: 0,
				Field: "test_field",
			},
			wantErr: false,
		},
		{
			name: "float value",
			validator: RequiredFieldValidator{
				Value: 3.14,
				Field: "test_field",
			},
			wantErr: false,
		},
		{
			name: "struct value",
			validator: RequiredFieldValidator{
				Value: struct{}{},
				Field: "test_field",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransactionTypeValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator TransactionTypeValidator
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty type (valid)",
			validator: TransactionTypeValidator{Type: ""},
			wantErr:   false,
		},
		{
			name:      "all type",
			validator: TransactionTypeValidator{Type: "all"},
			wantErr:   false,
		},
		{
			name:      "withdrawal type",
			validator: TransactionTypeValidator{Type: "withdrawal"},
			wantErr:   false,
		},
		{
			name:      "withdrawals type",
			validator: TransactionTypeValidator{Type: "withdrawals"},
			wantErr:   false,
		},
		{
			name:      "expense type",
			validator: TransactionTypeValidator{Type: "expense"},
			wantErr:   false,
		},
		{
			name:      "deposit type",
			validator: TransactionTypeValidator{Type: "deposit"},
			wantErr:   false,
		},
		{
			name:      "deposits type",
			validator: TransactionTypeValidator{Type: "deposits"},
			wantErr:   false,
		},
		{
			name:      "income type",
			validator: TransactionTypeValidator{Type: "income"},
			wantErr:   false,
		},
		{
			name:      "transfer type",
			validator: TransactionTypeValidator{Type: "transfer"},
			wantErr:   false,
		},
		{
			name:      "transfers type",
			validator: TransactionTypeValidator{Type: "transfers"},
			wantErr:   false,
		},
		{
			name:      "opening_balance type",
			validator: TransactionTypeValidator{Type: "opening_balance"},
			wantErr:   false,
		},
		{
			name:      "reconciliation type",
			validator: TransactionTypeValidator{Type: "reconciliation"},
			wantErr:   false,
		},
		{
			name:      "invalid type",
			validator: TransactionTypeValidator{Type: "invalid_type"},
			wantErr:   true,
			errMsg:    "validation error in field 'type': invalid transaction type: invalid_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompositeValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator CompositeValidator
		wantErr   bool
		errCount  int
	}{
		{
			name: "all validators pass",
			validator: CompositeValidator{
				Validators: []Validator{
					DateValidator{Value: "2024-01-15", Field: "date"},
					PaginationValidator{Page: 1, Limit: 10},
				},
			},
			wantErr: false,
		},
		{
			name: "first validator fails",
			validator: CompositeValidator{
				Validators: []Validator{
					DateValidator{Value: "invalid", Field: "date"},
					PaginationValidator{Page: 1, Limit: 10},
				},
			},
			wantErr: true,
		},
		{
			name: "second validator fails",
			validator: CompositeValidator{
				Validators: []Validator{
					DateValidator{Value: "2024-01-15", Field: "date"},
					PaginationValidator{Page: -1, Limit: 10},
				},
			},
			wantErr: true,
		},
		{
			name: "multiple validators fail",
			validator: CompositeValidator{
				Validators: []Validator{
					DateValidator{Value: "invalid", Field: "date"},
					PaginationValidator{Page: -1, Limit: -10},
				},
			},
			wantErr:  true,
			errCount: 2,
		},
		{
			name: "empty validators list",
			validator: CompositeValidator{
				Validators: []Validator{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Test ValidateAll
			errors := tt.validator.ValidateAll()
			if tt.errCount > 0 {
				assert.Len(t, errors, tt.errCount)
			} else if tt.wantErr {
				assert.NotEmpty(t, errors)
			} else {
				assert.Empty(t, errors)
			}
		})
	}
}

func TestAmountValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator AmountValidator
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid amount",
			validator: AmountValidator{
				Amount: "100.50",
				Field:  "amount",
			},
			wantErr: false,
		},
		{
			name: "empty amount when not required",
			validator: AmountValidator{
				Amount:   "",
				Field:    "amount",
				Required: false,
			},
			wantErr: false,
		},
		{
			name: "empty amount when required",
			validator: AmountValidator{
				Amount:   "",
				Field:    "amount",
				Required: true,
			},
			wantErr: true,
			errMsg:  "validation error in field 'amount': amount is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "test_field",
		Message: "test message",
	}

	expected := "validation error in field 'test_field': test message"
	assert.Equal(t, expected, err.Error())
}

func TestCombineErrors(t *testing.T) {
	tests := []struct {
		name    string
		errors  []error
		wantNil bool
		wantMsg string
	}{
		{
			name:    "no errors",
			errors:  []error{},
			wantNil: true,
		},
		{
			name:    "nil errors",
			errors:  nil,
			wantNil: true,
		},
		{
			name: "single error",
			errors: []error{
				ValidationError{Field: "field1", Message: "error1"},
			},
			wantNil: false,
			wantMsg: "validation error in field 'field1': error1",
		},
		{
			name: "multiple errors",
			errors: []error{
				ValidationError{Field: "field1", Message: "error1"},
				ValidationError{Field: "field2", Message: "error2"},
			},
			wantNil: false,
			wantMsg: "validation error in field 'field1': error1; validation error in field 'field2': error2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CombineErrors(tt.errors)
			if tt.wantNil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				if tt.wantMsg != "" {
					assert.Equal(t, tt.wantMsg, err.Error())
				}
			}
		})
	}
}