package dto

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

func TestSerializationFormats(t *testing.T) {
	// Test data
	account := &Account{
		Id:     "acc-123",
		Name:   "Test Account",
		Type:   "asset",
		Active: true,
		Notes:  ptrString("Some notes"),
	}
	
	budget := &Budget{
		Id:     "bud-456",
		Name:   "Monthly Budget",
		Active: true,
		Notes:  "Budget notes",
		Spent: Spent{
			Sum:          "100.50",
			CurrencyCode: "USD",
		},
	}
	
	transaction := &Transaction{
		Id:              "trans-789",
		Amount:          "250.00",
		Description:     "Test Transaction",
		Date:            time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		SourceId:        "acc-1",
		SourceName:      "Source Account",
		DestinationId:   "acc-2",
		DestinationName: "Dest Account",
		Type:            "withdrawal",
		CurrencyCode:    "EUR",
		Tags:            []string{"tag1", "tag2"},
	}
	
	tests := []struct {
		name   string
		format SerializationFormat
		dto    interface{}
		check  func(t *testing.T, data []byte)
	}{
		{
			name:   "Account JSON",
			format: FormatJSON,
			dto:    account,
			check: func(t *testing.T, data []byte) {
				var decoded Account
				if err := json.Unmarshal(data, &decoded); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				if decoded.Id != account.Id {
					t.Errorf("Id mismatch: got %v, want %v", decoded.Id, account.Id)
				}
				if decoded.Name != account.Name {
					t.Errorf("Name mismatch: got %v, want %v", decoded.Name, account.Name)
				}
			},
		},
		{
			name:   "Account YAML",
			format: FormatYAML,
			dto:    account,
			check: func(t *testing.T, data []byte) {
				var decoded Account
				if err := yaml.Unmarshal(data, &decoded); err != nil {
					t.Fatalf("Failed to unmarshal YAML: %v", err)
				}
				if decoded.Id != account.Id {
					t.Errorf("Id mismatch: got %v, want %v", decoded.Id, account.Id)
				}
				if decoded.Name != account.Name {
					t.Errorf("Name mismatch: got %v, want %v", decoded.Name, account.Name)
				}
			},
		},
		{
			name:   "Account MessagePack",
			format: FormatMessagePack,
			dto:    account,
			check: func(t *testing.T, data []byte) {
				var decoded Account
				if err := msgpack.Unmarshal(data, &decoded); err != nil {
					t.Fatalf("Failed to unmarshal MessagePack: %v", err)
				}
				if decoded.Id != account.Id {
					t.Errorf("Id mismatch: got %v, want %v", decoded.Id, account.Id)
				}
				if decoded.Name != account.Name {
					t.Errorf("Name mismatch: got %v, want %v", decoded.Name, account.Name)
				}
			},
		},
		{
			name:   "Budget JSON",
			format: FormatJSON,
			dto:    budget,
			check: func(t *testing.T, data []byte) {
				var decoded Budget
				if err := json.Unmarshal(data, &decoded); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				if decoded.Id != budget.Id {
					t.Errorf("Id mismatch: got %v, want %v", decoded.Id, budget.Id)
				}
				if decoded.Spent.Sum != budget.Spent.Sum {
					t.Errorf("Spent.Sum mismatch: got %v, want %v", decoded.Spent.Sum, budget.Spent.Sum)
				}
			},
		},
		{
			name:   "Transaction JSON with complex fields",
			format: FormatJSON,
			dto:    transaction,
			check: func(t *testing.T, data []byte) {
				var decoded Transaction
				if err := json.Unmarshal(data, &decoded); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				if decoded.Id != transaction.Id {
					t.Errorf("Id mismatch: got %v, want %v", decoded.Id, transaction.Id)
				}
				if len(decoded.Tags) != len(transaction.Tags) {
					t.Errorf("Tags length mismatch: got %v, want %v", len(decoded.Tags), len(transaction.Tags))
				}
				if !decoded.Date.Equal(transaction.Date) {
					t.Errorf("Date mismatch: got %v, want %v", decoded.Date, transaction.Date)
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serializer := DefaultSerializer
			
			// Marshal
			data, err := serializer.Marshal(tt.dto, tt.format)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			
			if len(data) == 0 {
				t.Error("Marshal returned empty data")
			}
			
			// Run format-specific checks
			tt.check(t, data)
		})
	}
}

func TestSerializableInterface(t *testing.T) {
	account := &Account{
		Id:     "acc-001",
		Name:   "Interface Test Account",
		Type:   "expense",
		Active: false,
	}
	
	tests := []struct {
		name   string
		format SerializationFormat
	}{
		{"JSON", FormatJSON},
		{"YAML", FormatYAML},
		{"MessagePack", FormatMessagePack},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Marshal method
			data, err := account.Marshal(tt.format)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			
			// Test Unmarshal method
			var decoded Account
			if err := decoded.Unmarshal(data, tt.format); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			
			// Verify
			if decoded.Id != account.Id {
				t.Errorf("Id mismatch: got %v, want %v", decoded.Id, account.Id)
			}
			if decoded.Name != account.Name {
				t.Errorf("Name mismatch: got %v, want %v", decoded.Name, account.Name)
			}
			if decoded.Type != account.Type {
				t.Errorf("Type mismatch: got %v, want %v", decoded.Type, account.Type)
			}
			if decoded.Active != account.Active {
				t.Errorf("Active mismatch: got %v, want %v", decoded.Active, account.Active)
			}
		})
	}
}

func TestSerializerWithWriter(t *testing.T) {
	budget := &Budget{
		Id:     "bud-001",
		Name:   "Writer Test Budget",
		Active: true,
		Spent: Spent{
			Sum:          "500.00",
			CurrencyCode: "GBP",
		},
	}
	
	tests := []struct {
		name   string
		format SerializationFormat
	}{
		{"JSON", FormatJSON},
		{"YAML", FormatYAML},
		{"MessagePack", FormatMessagePack},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test MarshalTo
			var buf bytes.Buffer
			if err := budget.MarshalTo(&buf, tt.format); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}
			
			if buf.Len() == 0 {
				t.Error("MarshalTo wrote no data")
			}
			
			// Test UnmarshalFrom
			var decoded Budget
			reader := bytes.NewReader(buf.Bytes())
			if err := decoded.UnmarshalFrom(reader, tt.format); err != nil {
				t.Fatalf("UnmarshalFrom failed: %v", err)
			}
			
			// Verify
			if decoded.Id != budget.Id {
				t.Errorf("Id mismatch: got %v, want %v", decoded.Id, budget.Id)
			}
			if decoded.Name != budget.Name {
				t.Errorf("Name mismatch: got %v, want %v", decoded.Name, budget.Name)
			}
		})
	}
}

func TestListSerialization(t *testing.T) {
	accountList := &AccountList{
		Data: []Account{
			{Id: "1", Name: "Account 1", Type: "asset", Active: true},
			{Id: "2", Name: "Account 2", Type: "expense", Active: false},
		},
		Pagination: Pagination{
			Count:       2,
			Total:       10,
			CurrentPage: 1,
			PerPage:     2,
			TotalPages:  5,
		},
	}
	
	tests := []struct {
		name   string
		format SerializationFormat
	}{
		{"JSON", FormatJSON},
		{"YAML", FormatYAML},
		{"MessagePack", FormatMessagePack},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := accountList.Marshal(tt.format)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			
			// Unmarshal
			var decoded AccountList
			if err := decoded.Unmarshal(data, tt.format); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			
			// Verify data
			if len(decoded.Data) != len(accountList.Data) {
				t.Errorf("Data length mismatch: got %v, want %v", len(decoded.Data), len(accountList.Data))
			}
			
			for i, account := range decoded.Data {
				if account.Id != accountList.Data[i].Id {
					t.Errorf("Account[%d].Id mismatch: got %v, want %v", i, account.Id, accountList.Data[i].Id)
				}
				if account.Name != accountList.Data[i].Name {
					t.Errorf("Account[%d].Name mismatch: got %v, want %v", i, account.Name, accountList.Data[i].Name)
				}
			}
			
			// Verify pagination
			if decoded.Pagination.Total != accountList.Pagination.Total {
				t.Errorf("Pagination.Total mismatch: got %v, want %v", decoded.Pagination.Total, accountList.Pagination.Total)
			}
			if decoded.Pagination.CurrentPage != accountList.Pagination.CurrentPage {
				t.Errorf("Pagination.CurrentPage mismatch: got %v, want %v", decoded.Pagination.CurrentPage, accountList.Pagination.CurrentPage)
			}
		})
	}
}

func TestSerializerOptions(t *testing.T) {
	account := &Account{
		Id:     "acc-opt",
		Name:   "Options Test",
		Type:   "asset",
		Active: true,
	}
	
	t.Run("JSON with indentation", func(t *testing.T) {
		serializer := &Serializer{
			JSONIndent:     "    ",
			JSONPrefix:     "",
			JSONEscapeHTML: false,
		}
		
		data, err := serializer.Marshal(account, FormatJSON)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		
		// Check for indentation
		if !strings.Contains(string(data), "\n    ") {
			t.Error("JSON output should contain indentation")
		}
	})
	
	t.Run("JSON compact", func(t *testing.T) {
		serializer := &Serializer{
			JSONIndent: "",
			JSONPrefix: "",
		}
		
		data, err := serializer.Marshal(account, FormatJSON)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		
		// Check for no newlines (compact)
		if strings.Contains(string(data), "\n") {
			t.Error("Compact JSON should not contain newlines")
		}
	})
	
	t.Run("YAML with custom indent", func(t *testing.T) {
		serializer := &Serializer{
			YAMLIndent: 4,
		}
		
		data, err := serializer.Marshal(account, FormatYAML)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		
		// YAML should be generated without error
		if len(data) == 0 {
			t.Error("YAML output should not be empty")
		}
	})
}

func TestUnsupportedFormat(t *testing.T) {
	account := &Account{
		Id:   "acc-1",
		Name: "Test",
	}
	
	serializer := DefaultSerializer
	
	// Test unsupported format
	_, err := serializer.Marshal(account, "unsupported")
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported serialization format") {
		t.Errorf("Error message should mention unsupported format: %v", err)
	}
	
	// Test unmarshal with unsupported format
	err = serializer.Unmarshal([]byte("{}"), account, "unsupported")
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}

func TestConvenienceFunctions(t *testing.T) {
	category := &Category{
		Id:    "cat-123",
		Name:  "Test Category",
		Notes: "Some notes",
	}
	
	t.Run("JSON convenience functions", func(t *testing.T) {
		// Marshal
		data, err := MarshalJSON(category)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}
		
		// Unmarshal
		var decoded Category
		if err := UnmarshalJSON(data, &decoded); err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}
		
		if decoded.Id != category.Id {
			t.Errorf("Id mismatch: got %v, want %v", decoded.Id, category.Id)
		}
	})
	
	t.Run("YAML convenience functions", func(t *testing.T) {
		// Marshal
		data, err := MarshalYAML(category)
		if err != nil {
			t.Fatalf("MarshalYAML failed: %v", err)
		}
		
		// Unmarshal
		var decoded Category
		if err := UnmarshalYAML(data, &decoded); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}
		
		if decoded.Id != category.Id {
			t.Errorf("Id mismatch: got %v, want %v", decoded.Id, category.Id)
		}
	})
	
	t.Run("MessagePack convenience functions", func(t *testing.T) {
		// Marshal
		data, err := MarshalMessagePack(category)
		if err != nil {
			t.Fatalf("MarshalMessagePack failed: %v", err)
		}
		
		// Unmarshal
		var decoded Category
		if err := UnmarshalMessagePack(data, &decoded); err != nil {
			t.Fatalf("UnmarshalMessagePack failed: %v", err)
		}
		
		if decoded.Id != category.Id {
			t.Errorf("Id mismatch: got %v, want %v", decoded.Id, category.Id)
		}
	})
}

func TestComplexDataStructures(t *testing.T) {
	transactionGroup := &TransactionGroup{
		Id:         "tg-001",
		GroupTitle: "Complex Transaction Group",
		Transactions: []Transaction{
			{
				Id:              "t1",
				Amount:          "100.00",
				Description:     "Transaction 1",
				Date:            time.Now(),
				SourceId:        "s1",
				SourceName:      "Source 1",
				DestinationId:   "d1",
				DestinationName: "Dest 1",
				Type:            "withdrawal",
				CurrencyCode:    "USD",
				Tags:            []string{"tag1", "tag2"},
				CategoryId:      ptrString("cat-1"),
				CategoryName:    ptrString("Category 1"),
				BudgetId:        ptrString("bud-1"),
				BudgetName:      ptrString("Budget 1"),
				Notes:           ptrString("Some notes"),
			},
			{
				Id:              "t2",
				Amount:          "200.00",
				Description:     "Transaction 2",
				Date:            time.Now().Add(24 * time.Hour),
				SourceId:        "s2",
				SourceName:      "Source 2",
				DestinationId:   "d2",
				DestinationName: "Dest 2",
				Type:            "deposit",
				CurrencyCode:    "EUR",
				Tags:            []string{},
			},
		},
	}
	
	formats := []SerializationFormat{FormatJSON, FormatYAML, FormatMessagePack}
	
	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			// Marshal
			data, err := transactionGroup.Marshal(format)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			
			// Unmarshal
			var decoded TransactionGroup
			if err := decoded.Unmarshal(data, format); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			
			// Verify
			if decoded.Id != transactionGroup.Id {
				t.Errorf("Id mismatch: got %v, want %v", decoded.Id, transactionGroup.Id)
			}
			
			if len(decoded.Transactions) != len(transactionGroup.Transactions) {
				t.Fatalf("Transactions count mismatch: got %v, want %v", 
					len(decoded.Transactions), len(transactionGroup.Transactions))
			}
			
			// Check first transaction details
			tx1 := decoded.Transactions[0]
			origTx1 := transactionGroup.Transactions[0]
			
			if tx1.Id != origTx1.Id {
				t.Errorf("Transaction[0].Id mismatch: got %v, want %v", tx1.Id, origTx1.Id)
			}
			if tx1.Amount != origTx1.Amount {
				t.Errorf("Transaction[0].Amount mismatch: got %v, want %v", tx1.Amount, origTx1.Amount)
			}
			if len(tx1.Tags) != len(origTx1.Tags) {
				t.Errorf("Transaction[0].Tags length mismatch: got %v, want %v", 
					len(tx1.Tags), len(origTx1.Tags))
			}
			
			// Check optional fields
			if (tx1.CategoryId == nil) != (origTx1.CategoryId == nil) {
				t.Error("Transaction[0].CategoryId nil mismatch")
			}
			if tx1.CategoryId != nil && *tx1.CategoryId != *origTx1.CategoryId {
				t.Errorf("Transaction[0].CategoryId mismatch: got %v, want %v", 
					*tx1.CategoryId, *origTx1.CategoryId)
			}
		})
	}
}

func TestNilHandling(t *testing.T) {
	account := &Account{
		Id:     "acc-nil",
		Name:   "Nil Test",
		Type:   "asset",
		Active: true,
		Notes:  nil, // Explicitly nil
	}
	
	formats := []SerializationFormat{FormatJSON, FormatYAML, FormatMessagePack}
	
	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			// Marshal
			data, err := account.Marshal(format)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			
			// Unmarshal
			var decoded Account
			if err := decoded.Unmarshal(data, format); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			
			// Verify nil field remains nil
			if decoded.Notes != nil {
				t.Errorf("Notes should be nil but got: %v", decoded.Notes)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that data survives a round trip through all formats
	budgetList := &BudgetList{
		Data: []Budget{
			{
				Id:     "b1",
				Name:   "Budget 1",
				Active: true,
				Notes:  "Note 1",
				Spent: Spent{
					Sum:          "100.00",
					CurrencyCode: "USD",
				},
			},
			{
				Id:     "b2",
				Name:   "Budget 2",
				Active: false,
				Notes:  nil,
				Spent: Spent{
					Sum:          "200.00",
					CurrencyCode: "EUR",
				},
			},
		},
		Pagination: Pagination{
			Count:       2,
			Total:       2,
			CurrentPage: 1,
			PerPage:     10,
			TotalPages:  1,
		},
	}
	
	formats := []SerializationFormat{FormatJSON, FormatYAML, FormatMessagePack}
	
	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			// First round trip
			data1, err := budgetList.Marshal(format)
			if err != nil {
				t.Fatalf("First marshal failed: %v", err)
			}
			
			var decoded1 BudgetList
			if err := decoded1.Unmarshal(data1, format); err != nil {
				t.Fatalf("First unmarshal failed: %v", err)
			}
			
			// Second round trip
			data2, err := decoded1.Marshal(format)
			if err != nil {
				t.Fatalf("Second marshal failed: %v", err)
			}
			
			var decoded2 BudgetList
			if err := decoded2.Unmarshal(data2, format); err != nil {
				t.Fatalf("Second unmarshal failed: %v", err)
			}
			
			// Data should be equivalent after round trips
			if !reflect.DeepEqual(decoded1.Data, decoded2.Data) {
				t.Error("Data changed after round trip")
			}
			
			// For JSON and YAML, the serialized data should be the same
			if format != FormatMessagePack {
				if !bytes.Equal(data1, data2) {
					t.Errorf("Serialized data changed after round trip for %s", format)
				}
			}
		})
	}
}