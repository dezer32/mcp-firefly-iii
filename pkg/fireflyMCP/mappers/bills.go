package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapBillArrayToBillList converts client.BillArray to BillList DTO
func MapBillArrayToBillList(billArray *client.BillArray) *dto.BillList {
	if billArray == nil {
		return nil
	}

	return MapArrayToList(
		billArray,
		func(billRead client.BillRead) dto.Bill {
			billPtr := MapBillReadToBill(&billRead)
			if billPtr != nil {
				return *billPtr
			}
			return dto.NewBillBuilder().Build()
		},
		func() *dto.BillList { return &dto.BillList{} },
	)
}

// MapBillReadToBill converts client.BillRead to Bill DTO
func MapBillReadToBill(billRead *client.BillRead) *dto.Bill {
	if billRead == nil {
		return nil
	}

	// Start building the Bill
	builder := dto.NewBillBuilder().
		WithId(billRead.Id).
		WithActive(GetBoolValue(billRead.Attributes.Active)).
		WithName(billRead.Attributes.Name).
		WithAmountMin(billRead.Attributes.AmountMin).
		WithAmountMax(billRead.Attributes.AmountMax).
		WithDate(billRead.Attributes.Date).
		WithRepeatFreq(string(billRead.Attributes.RepeatFreq)).
		WithCurrencyCode(GetStringValue(billRead.Attributes.CurrencyCode)).
		WithNotes(billRead.Attributes.Notes)

	// Handle skip
	skip := 0
	if billRead.Attributes.Skip != nil {
		skip = int(*billRead.Attributes.Skip)
	}
	builder = builder.WithSkip(skip)

	// Handle next expected match
	if billRead.Attributes.NextExpectedMatch != nil {
		builder = builder.WithNextExpectedMatch(billRead.Attributes.NextExpectedMatch)
	}

	// Handle paid dates
	paidDates := []dto.PaidDate{}
	if billRead.Attributes.PaidDates != nil {
		for _, pd := range *billRead.Attributes.PaidDates {
			paidDateBuilder := dto.NewPaidDateBuilder().
				WithTransactionGroupId(pd.TransactionGroupId).
				WithTransactionJournalId(pd.TransactionJournalId)
			
			if pd.Date != nil {
				paidDateBuilder = paidDateBuilder.WithDate(pd.Date)
			}
			
			paidDates = append(paidDates, paidDateBuilder.Build())
		}
	}
	builder = builder.WithPaidDates(paidDates)

	bill := builder.Build()
	return &bill
}