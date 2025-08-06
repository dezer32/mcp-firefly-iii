package mappers

import (
	"time"
	
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
			return dto.Bill{}
		},
		func() *dto.BillList { return &dto.BillList{} },
	)
}

// MapBillReadToBill converts client.BillRead to Bill DTO
func MapBillReadToBill(billRead *client.BillRead) *dto.Bill {
	if billRead == nil {
		return nil
	}

	bill := &dto.Bill{
		Id:           billRead.Id,
		Active:       GetBoolValue(billRead.Attributes.Active),
		Name:         billRead.Attributes.Name,
		AmountMin:    billRead.Attributes.AmountMin,
		AmountMax:    billRead.Attributes.AmountMax,
		Date:         time.Time{},
		RepeatFreq:   string(billRead.Attributes.RepeatFreq),
		Skip:         0,
		CurrencyCode: GetStringValue(billRead.Attributes.CurrencyCode),
		Notes:        billRead.Attributes.Notes,
		PaidDates:    []dto.PaidDate{},
	}

	// Handle date - Date is already a time.Time, not a pointer
	bill.Date = billRead.Attributes.Date

	// Handle skip
	if billRead.Attributes.Skip != nil {
		bill.Skip = int(*billRead.Attributes.Skip)
	}

	// Handle next expected match
	if billRead.Attributes.NextExpectedMatch != nil {
		bill.NextExpectedMatch = billRead.Attributes.NextExpectedMatch
	}

	// Handle paid dates
	if billRead.Attributes.PaidDates != nil {
		for _, pd := range *billRead.Attributes.PaidDates {
			paidDate := dto.PaidDate{
				Date:                 nil,
				TransactionGroupId:   pd.TransactionGroupId,
				TransactionJournalId: pd.TransactionJournalId,
			}
			if pd.Date != nil {
				paidDate.Date = pd.Date
			}
			bill.PaidDates = append(bill.PaidDates, paidDate)
		}
	}

	return bill
}