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

	transactionList := &dto.TransactionList{
		Data: []dto.TransactionGroup{},
	}

	// Map transaction data - each TransactionRead becomes a TransactionGroup
	for _, transactionRead := range transactionArray.Data {
		group := dto.TransactionGroup{
			Id:           transactionRead.Id,
			GroupTitle:   GetStringValue(transactionRead.Attributes.GroupTitle),
			Transactions: []dto.Transaction{},
		}

		for _, split := range transactionRead.Attributes.Transactions {
			transaction := mapTransactionSplitToTransaction(&split)
			if transaction != nil {
				group.Transactions = append(group.Transactions, *transaction)
			}
		}

		transactionList.Data = append(transactionList.Data, group)
	}

	// Map pagination
	if transactionArray.Meta.Pagination != nil {
		transactionList.Pagination = mapPaginationToDTO(transactionArray.Meta.Pagination)
	}

	return transactionList
}

// MapTransactionReadToTransactionGroup converts a TransactionRead to TransactionGroup DTO.
// Maps transaction splits and metadata to the simplified group format.
// Returns nil if the input is nil
func MapTransactionReadToTransactionGroup(transactionRead *client.TransactionRead) *dto.TransactionGroup {
	if transactionRead == nil {
		return nil
	}

	group := &dto.TransactionGroup{
		Id:           transactionRead.Id,
		GroupTitle:   GetStringValue(transactionRead.Attributes.GroupTitle),
		Transactions: make([]dto.Transaction, 0),
	}

	// Map all transaction splits
	for _, split := range transactionRead.Attributes.Transactions {
		if txn := mapTransactionSplitToTransaction(&split); txn != nil {
			group.Transactions = append(group.Transactions, *txn)
		}
	}

	return group
}

// mapTransactionSplitToTransaction converts a TransactionSplit to Transaction DTO
func mapTransactionSplitToTransaction(split *client.TransactionSplit) *dto.Transaction {
	if split == nil {
		return nil
	}

	transaction := &dto.Transaction{
		Id:              GetStringValue(split.TransactionJournalId),
		Amount:          split.Amount,
		Date:            split.Date,
		Description:     split.Description,
		SourceId:        GetStringValue(split.SourceId),
		SourceName:      GetStringValue(split.SourceName),
		DestinationId:   GetStringValue(split.DestinationId),
		DestinationName: GetStringValue(split.DestinationName),
		Type:            string(split.Type),
		CategoryId:      split.CategoryId,
		CategoryName:    split.CategoryName,
		BudgetId:        split.BudgetId,
		BudgetName:      split.BudgetName,
		Notes:           split.Notes,
		Reconciled:      GetBoolValue(split.Reconciled),
		CurrencyCode:    GetStringValue(split.CurrencyCode),
		DestinationType: GetAccountTypeString(split.DestinationType),
	}

	// Map tags if present
	if split.Tags != nil && len(*split.Tags) > 0 {
		transaction.Tags = *split.Tags
	} else {
		transaction.Tags = []string{}
	}

	// Handle bill fields
	if split.BillId != nil {
		transaction.BillId = split.BillId
	}
	if split.BillName != nil {
		transaction.BillName = split.BillName
	}

	return transaction
}

// MapTransactionArrayToTransactionGroupList is an alias for MapTransactionArrayToTransactionList
func MapTransactionArrayToTransactionGroupList(transactionArray *client.TransactionArray) *dto.TransactionList {
	return MapTransactionArrayToTransactionList(transactionArray)
}

// mapPaginationToDTO converts client pagination to Pagination DTO
func mapPaginationToDTO(pagination interface{}) dto.Pagination {
	// Return empty pagination if nil
	if pagination == nil {
		return dto.Pagination{}
	}
	
	// Try to cast to the correct type - note the actual type in client may vary
	// We'll need to check what the actual type is
	return dto.Pagination{
		// These fields will be populated based on the actual pagination structure
		// For now, return empty pagination
	}
}