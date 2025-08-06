package builders

import (
	"testing"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAccountParamsBuilder(t *testing.T) {
	tests := []struct {
		name      string
		build     func() *ListAccountParamsBuilder
		wantError bool
		errorMsg  string
		validate  func(t *testing.T, params interface{})
	}{
		{
			name: "with asset accounts filter",
			build: func() *ListAccountParamsBuilder {
				b := NewListAccountParamsBuilder()
				b.WithAssetAccounts()
				b.WithLimit(20)
				b.WithPage(1)
				return b
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.ListAccountParams)
				require.True(t, ok)
				require.NotNil(t, p.Type)
				assert.Equal(t, "asset", string(*p.Type))
				require.NotNil(t, p.Limit)
				assert.Equal(t, int32(20), *p.Limit)
				require.NotNil(t, p.Page)
				assert.Equal(t, int32(1), *p.Page)
			},
		},
		{
			name: "with expense accounts filter",
			build: func() *ListAccountParamsBuilder {
				return NewListAccountParamsBuilder().
					WithExpenseAccounts()
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.ListAccountParams)
				require.True(t, ok)
				require.NotNil(t, p.Type)
				assert.Equal(t, "expense", string(*p.Type))
			},
		},
		{
			name: "with revenue accounts filter",
			build: func() *ListAccountParamsBuilder {
				return NewListAccountParamsBuilder().
					WithRevenueAccounts()
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.ListAccountParams)
				require.True(t, ok)
				require.NotNil(t, p.Type)
				assert.Equal(t, "revenue", string(*p.Type))
			},
		},
		{
			name: "with custom type",
			build: func() *ListAccountParamsBuilder {
				return NewListAccountParamsBuilder().
					WithType("liability")
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.ListAccountParams)
				require.True(t, ok)
				require.NotNil(t, p.Type)
				assert.Equal(t, "liability", string(*p.Type))
			},
		},
		{
			name: "invalid account type",
			build: func() *ListAccountParamsBuilder {
				return NewListAccountParamsBuilder().
					WithType("invalid_type")
			},
			wantError: true,
			errorMsg:  "invalid account type",
		},
		{
			name: "no filters",
			build: func() *ListAccountParamsBuilder {
				return NewListAccountParamsBuilder()
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.ListAccountParams)
				require.True(t, ok)
				assert.Nil(t, p.Type)
				assert.Nil(t, p.Limit)
				assert.Nil(t, p.Page)
			},
		},
		{
			name: "pagination only",
			build: func() *ListAccountParamsBuilder {
				b := NewListAccountParamsBuilder()
				b.WithLimit(50)
				b.WithPage(3)
				return b
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.ListAccountParams)
				require.True(t, ok)
				assert.Nil(t, p.Type)
				require.NotNil(t, p.Limit)
				assert.Equal(t, int32(50), *p.Limit)
				require.NotNil(t, p.Page)
				assert.Equal(t, int32(3), *p.Page)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.build()
			params, err := builder.Build()

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, params)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, params)
				if tt.validate != nil {
					tt.validate(t, params)
				}
			}
		})
	}
}

func TestSearchAccountParamsBuilder(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		build     func(b *SearchAccountParamsBuilder) *SearchAccountParamsBuilder
		wantError bool
		errorMsg  string
		validate  func(t *testing.T, params interface{})
	}{
		{
			name:  "basic search",
			query: "savings",
			build: func(b *SearchAccountParamsBuilder) *SearchAccountParamsBuilder {
				return b
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.SearchAccountsParams)
				require.True(t, ok)
				assert.Equal(t, "savings", p.Query)
				assert.Equal(t, "all", string(p.Field))
			},
		},
		{
			name:  "search in name field",
			query: "checking",
			build: func(b *SearchAccountParamsBuilder) *SearchAccountParamsBuilder {
				b.WithField("name")
				return b
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.SearchAccountsParams)
				require.True(t, ok)
				assert.Equal(t, "checking", p.Query)
				assert.Equal(t, "name", string(p.Field))
			},
		},
		{
			name:  "search with pagination",
			query: "account",
			build: func(b *SearchAccountParamsBuilder) *SearchAccountParamsBuilder {
				b.WithLimit(10)
				b.WithPage(2)
				return b
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.SearchAccountsParams)
				require.True(t, ok)
				assert.Equal(t, "account", p.Query)
				require.NotNil(t, p.Limit)
				assert.Equal(t, int32(10), *p.Limit)
				require.NotNil(t, p.Page)
				assert.Equal(t, int32(2), *p.Page)
			},
		},
		{
			name:  "search in IBAN field",
			query: "DE89",
			build: func(b *SearchAccountParamsBuilder) *SearchAccountParamsBuilder {
				b.WithField("iban")
				return b
			},
			wantError: false,
			validate: func(t *testing.T, params interface{}) {
				p, ok := params.(*client.SearchAccountsParams)
				require.True(t, ok)
				assert.Equal(t, "DE89", p.Query)
				assert.Equal(t, "iban", string(p.Field))
			},
		},
		{
			name:  "empty query",
			query: "",
			build: func(b *SearchAccountParamsBuilder) *SearchAccountParamsBuilder {
				return b
			},
			wantError: true,
			errorMsg:  "search query cannot be empty",
		},
		{
			name:  "invalid field",
			query: "test",
			build: func(b *SearchAccountParamsBuilder) *SearchAccountParamsBuilder {
				b.WithField("invalid")
				return b
			},
			wantError: true,
			errorMsg:  "invalid search field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchAccountParamsBuilder(tt.query)
			builder = tt.build(builder)
			params, err := builder.Build()

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, params)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, params)
				if tt.validate != nil {
					tt.validate(t, params)
				}
			}
		})
	}
}