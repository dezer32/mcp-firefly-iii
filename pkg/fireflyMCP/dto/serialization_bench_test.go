package dto

import (
	"testing"
	"time"
)

// Sample data for benchmarks
var (
	benchAccount = &Account{
		Id:     "acc-bench-123",
		Name:   "Benchmark Test Account",
		Type:   "asset",
		Active: true,
		Notes:  ptrString("This is a test account for benchmarking serialization performance"),
	}
	
	benchTransaction = &Transaction{
		Id:              "trans-bench-456",
		Amount:          "1234.56",
		Description:     "Benchmark transaction for testing serialization performance",
		Date:            time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		SourceId:        "src-account-1",
		SourceName:      "Source Account Name",
		DestinationId:   "dest-account-2",
		DestinationName: "Destination Account Name",
		Type:            "withdrawal",
		CurrencyCode:    "USD",
		DestinationType: "expense",
		Tags:            []string{"benchmark", "test", "performance"},
		CategoryId:      ptrString("cat-123"),
		CategoryName:    ptrString("Test Category"),
		BudgetId:        ptrString("bud-456"),
		BudgetName:      ptrString("Test Budget"),
		Notes:           ptrString("Additional notes for the transaction"),
		Reconciled:      true,
	}
	
	benchAccountList = &AccountList{
		Data: []Account{
			*benchAccount,
			{Id: "acc-2", Name: "Account 2", Type: "expense", Active: false},
			{Id: "acc-3", Name: "Account 3", Type: "revenue", Active: true},
			{Id: "acc-4", Name: "Account 4", Type: "asset", Active: true},
			{Id: "acc-5", Name: "Account 5", Type: "liability", Active: false},
		},
		Pagination: Pagination{
			Count:       5,
			Total:       100,
			CurrentPage: 1,
			PerPage:     5,
			TotalPages:  20,
		},
	}
	
	benchTransactionGroup = &TransactionGroup{
		Id:         "tg-bench-789",
		GroupTitle: "Benchmark Transaction Group",
		Transactions: []Transaction{
			*benchTransaction,
			{
				Id:              "trans-2",
				Amount:          "500.00",
				Description:     "Second transaction",
				Date:            time.Now(),
				SourceId:        "src-2",
				SourceName:      "Source 2",
				DestinationId:   "dest-2",
				DestinationName: "Dest 2",
				Type:            "deposit",
				CurrencyCode:    "EUR",
				Tags:            []string{"tag1", "tag2"},
			},
			{
				Id:              "trans-3",
				Amount:          "750.25",
				Description:     "Third transaction",
				Date:            time.Now().Add(-24 * time.Hour),
				SourceId:        "src-3",
				SourceName:      "Source 3",
				DestinationId:   "dest-3",
				DestinationName: "Dest 3",
				Type:            "transfer",
				CurrencyCode:    "GBP",
				Tags:            []string{},
			},
		},
	}
)

// Benchmark JSON serialization

func BenchmarkMarshalJSON_Account(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchAccount.Marshal(FormatJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJSON_Account(b *testing.B) {
	data, _ := benchAccount.Marshal(FormatJSON)
	var account Account
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := account.Unmarshal(data, FormatJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalJSON_Transaction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchTransaction.Marshal(FormatJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJSON_Transaction(b *testing.B) {
	data, _ := benchTransaction.Marshal(FormatJSON)
	var transaction Transaction
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := transaction.Unmarshal(data, FormatJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalJSON_AccountList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchAccountList.Marshal(FormatJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJSON_AccountList(b *testing.B) {
	data, _ := benchAccountList.Marshal(FormatJSON)
	var list AccountList
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := list.Unmarshal(data, FormatJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalJSON_TransactionGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchTransactionGroup.Marshal(FormatJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark YAML serialization

func BenchmarkMarshalYAML_Account(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchAccount.Marshal(FormatYAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalYAML_Account(b *testing.B) {
	data, _ := benchAccount.Marshal(FormatYAML)
	var account Account
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := account.Unmarshal(data, FormatYAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalYAML_Transaction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchTransaction.Marshal(FormatYAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalYAML_Transaction(b *testing.B) {
	data, _ := benchTransaction.Marshal(FormatYAML)
	var transaction Transaction
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := transaction.Unmarshal(data, FormatYAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalYAML_AccountList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchAccountList.Marshal(FormatYAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalYAML_AccountList(b *testing.B) {
	data, _ := benchAccountList.Marshal(FormatYAML)
	var list AccountList
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := list.Unmarshal(data, FormatYAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark MessagePack serialization

func BenchmarkMarshalMessagePack_Account(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchAccount.Marshal(FormatMessagePack)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalMessagePack_Account(b *testing.B) {
	data, _ := benchAccount.Marshal(FormatMessagePack)
	var account Account
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := account.Unmarshal(data, FormatMessagePack)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalMessagePack_Transaction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchTransaction.Marshal(FormatMessagePack)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalMessagePack_Transaction(b *testing.B) {
	data, _ := benchTransaction.Marshal(FormatMessagePack)
	var transaction Transaction
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := transaction.Unmarshal(data, FormatMessagePack)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalMessagePack_AccountList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchAccountList.Marshal(FormatMessagePack)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalMessagePack_AccountList(b *testing.B) {
	data, _ := benchAccountList.Marshal(FormatMessagePack)
	var list AccountList
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := list.Unmarshal(data, FormatMessagePack)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalMessagePack_TransactionGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := benchTransactionGroup.Marshal(FormatMessagePack)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Comparative benchmarks - same data, different formats

func BenchmarkComparison_Marshal_SmallData(b *testing.B) {
	formats := []struct {
		name   string
		format SerializationFormat
	}{
		{"JSON", FormatJSON},
		{"YAML", FormatYAML},
		{"MessagePack", FormatMessagePack},
	}
	
	for _, f := range formats {
		b.Run(f.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := benchAccount.Marshal(f.format)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkComparison_Marshal_LargeData(b *testing.B) {
	formats := []struct {
		name   string
		format SerializationFormat
	}{
		{"JSON", FormatJSON},
		{"YAML", FormatYAML},
		{"MessagePack", FormatMessagePack},
	}
	
	for _, f := range formats {
		b.Run(f.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := benchTransactionGroup.Marshal(f.format)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkComparison_Unmarshal_SmallData(b *testing.B) {
	dataJSON, _ := benchAccount.Marshal(FormatJSON)
	dataYAML, _ := benchAccount.Marshal(FormatYAML)
	dataMsgPack, _ := benchAccount.Marshal(FormatMessagePack)
	
	formats := []struct {
		name   string
		format SerializationFormat
		data   []byte
	}{
		{"JSON", FormatJSON, dataJSON},
		{"YAML", FormatYAML, dataYAML},
		{"MessagePack", FormatMessagePack, dataMsgPack},
	}
	
	for _, f := range formats {
		b.Run(f.name, func(b *testing.B) {
			var account Account
			for i := 0; i < b.N; i++ {
				err := account.Unmarshal(f.data, f.format)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkComparison_Unmarshal_LargeData(b *testing.B) {
	dataJSON, _ := benchTransactionGroup.Marshal(FormatJSON)
	dataYAML, _ := benchTransactionGroup.Marshal(FormatYAML)
	dataMsgPack, _ := benchTransactionGroup.Marshal(FormatMessagePack)
	
	formats := []struct {
		name   string
		format SerializationFormat
		data   []byte
	}{
		{"JSON", FormatJSON, dataJSON},
		{"YAML", FormatYAML, dataYAML},
		{"MessagePack", FormatMessagePack, dataMsgPack},
	}
	
	for _, f := range formats {
		b.Run(f.name, func(b *testing.B) {
			var group TransactionGroup
			for i := 0; i < b.N; i++ {
				err := group.Unmarshal(f.data, f.format)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Memory allocation benchmarks

func BenchmarkMemoryAllocation_JSON(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data, _ := benchAccountList.Marshal(FormatJSON)
		var list AccountList
		_ = list.Unmarshal(data, FormatJSON)
	}
}

func BenchmarkMemoryAllocation_YAML(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data, _ := benchAccountList.Marshal(FormatYAML)
		var list AccountList
		_ = list.Unmarshal(data, FormatYAML)
	}
}

func BenchmarkMemoryAllocation_MessagePack(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data, _ := benchAccountList.Marshal(FormatMessagePack)
		var list AccountList
		_ = list.Unmarshal(data, FormatMessagePack)
	}
}

// Data size comparison
func BenchmarkDataSize_Comparison(b *testing.B) {
	// This isn't a performance benchmark but shows data size differences
	dataJSON, _ := benchTransactionGroup.Marshal(FormatJSON)
	dataYAML, _ := benchTransactionGroup.Marshal(FormatYAML)
	dataMsgPack, _ := benchTransactionGroup.Marshal(FormatMessagePack)
	
	b.Logf("Data sizes for TransactionGroup:")
	b.Logf("  JSON:        %d bytes", len(dataJSON))
	b.Logf("  YAML:        %d bytes", len(dataYAML))
	b.Logf("  MessagePack: %d bytes", len(dataMsgPack))
	
	// Calculate compression ratios
	if len(dataMsgPack) > 0 {
		b.Logf("Compression ratios (compared to JSON):")
		b.Logf("  YAML:        %.2f%%", float64(len(dataYAML))*100/float64(len(dataJSON)))
		b.Logf("  MessagePack: %.2f%%", float64(len(dataMsgPack))*100/float64(len(dataJSON)))
	}
}