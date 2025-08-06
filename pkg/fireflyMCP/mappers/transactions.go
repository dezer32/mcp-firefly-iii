package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapTransactionArrayToTransactionList converts client.TransactionArray to TransactionList DTO.
// Each TransactionRead in the array becomes a TransactionGroup with its transaction splits.
// Returns nil if the input is nil
func MapTransactionArrayToTransactionList(transactionArray *client.TransactionArray) *dto.TransactionList {
	if transactionArray == nil {
		return nil
	}

	listBuilder := dto.NewTransactionListBuilder()

	// Map transaction data - each TransactionRead becomes a TransactionGroup
	for _, transactionRead := range transactionArray.Data {
		groupBuilder := dto.NewTransactionGroupBuilder().
			WithId(transactionRead.Id).
			WithGroupTitle(GetStringValue(transactionRead.Attributes.GroupTitle))

		transactions := []dto.Transaction{}
		for _, split := range transactionRead.Attributes.Transactions {
			transaction := mapTransactionSplitToTransaction(&split)
			if transaction != nil {
				transactions = append(transactions, *transaction)
			}
		}
		groupBuilder = groupBuilder.WithTransactions(transactions)

		listBuilder = listBuilder.AddTransactionGroup(groupBuilder.Build())
	}

	// Map pagination
	if transactionArray.Meta.Pagination != nil {
		listBuilder = listBuilder.WithPagination(mapPaginationToDTO(transactionArray.Meta.Pagination))
	}

	transactionList := listBuilder.Build()
	return &transactionList
}

// MapTransactionReadToTransactionGroup converts a TransactionRead to TransactionGroup DTO.
// Maps transaction splits and metadata to the simplified group format.
// Returns nil if the input is nil
func MapTransactionReadToTransactionGroup(transactionRead *client.TransactionRead) *dto.TransactionGroup {
	if transactionRead == nil {
		return nil
	}

	groupBuilder := dto.NewTransactionGroupBuilder().
		WithId(transactionRead.Id).
		WithGroupTitle(GetStringValue(transactionRead.Attributes.GroupTitle))

	// Map all transaction splits
	transactions := []dto.Transaction{}
	for _, split := range transactionRead.Attributes.Transactions {
		if txn := mapTransactionSplitToTransaction(&split); txn != nil {
			transactions = append(transactions, *txn)
		}
	}
	groupBuilder = groupBuilder.WithTransactions(transactions)

	group := groupBuilder.Build()
	return &group
}

// mapTransactionSplitToTransaction converts a TransactionSplit to Transaction DTO
func mapTransactionSplitToTransaction(split *client.TransactionSplit) *dto.Transaction {
	if split == nil {
		return nil
	}

	builder := dto.NewTransactionBuilder().
		WithId(GetStringValue(split.TransactionJournalId)).
		WithAmount(split.Amount).
		WithDate(split.Date).
		WithDescription(split.Description).
		WithSourceId(GetStringValue(split.SourceId)).
		WithSourceName(GetStringValue(split.SourceName)).
		WithDestinationId(GetStringValue(split.DestinationId)).
		WithDestinationName(GetStringValue(split.DestinationName)).
		WithType(string(split.Type)).
		WithCategoryId(split.CategoryId).
		WithCategoryName(split.CategoryName).
		WithBudgetId(split.BudgetId).
		WithBudgetName(split.BudgetName).
		WithNotes(split.Notes).
		WithReconciled(GetBoolValue(split.Reconciled)).
		WithCurrencyCode(GetStringValue(split.CurrencyCode)).
		WithDestinationType(GetAccountTypeString(split.DestinationType))

	// Map tags if present
	tags := []string{}
	if split.Tags != nil && len(*split.Tags) > 0 {
		tags = *split.Tags
	}
	builder = builder.WithTags(tags)

	// Handle bill fields
	if split.BillId != nil {
		builder = builder.WithBillId(split.BillId)
	}
	if split.BillName != nil {
		builder = builder.WithBillName(split.BillName)
	}

	transaction := builder.Build()
	return &transaction
}

// MapTransactionArrayToTransactionGroupList is an alias for MapTransactionArrayToTransactionList
func MapTransactionArrayToTransactionGroupList(transactionArray *client.TransactionArray) *dto.TransactionList {
	return MapTransactionArrayToTransactionList(transactionArray)
}

// mapPaginationToDTO converts client pagination to Pagination DTO
func mapPaginationToDTO(pagination interface{}) dto.Pagination {
	// Use the generic mapper from generic.go
	return MapPaginationToDTO(pagination)
}