package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapBudgetArrayToBudgetList converts client.BudgetArray to BudgetList DTO.
// Maps budget data including spent amounts and pagination information.
// Returns nil if the input is nil
func MapBudgetArrayToBudgetList(budgetArray *client.BudgetArray) *dto.BudgetList {
	if budgetArray == nil {
		return nil
	}

	return MapArrayToList(
		budgetArray,
		func(budgetRead client.BudgetRead) dto.Budget {
			budget := dto.Budget{
				Id:     budgetRead.Id,
				Active: GetBoolValue(budgetRead.Attributes.Active),
				Name:   budgetRead.Attributes.Name,
				Notes:  budgetRead.Attributes.Notes,
				Spent: dto.Spent{
					Sum:          "0",
					CurrencyCode: "",
				},
			}

			// Map spent information if available
			if budgetRead.Attributes.Spent != nil && len(*budgetRead.Attributes.Spent) > 0 {
				for _, spent := range *budgetRead.Attributes.Spent {
					if spent.Sum != nil {
						budget.Spent.Sum = *spent.Sum
					}
					if spent.CurrencyCode != nil {
						budget.Spent.CurrencyCode = *spent.CurrencyCode
					}
					break // Take first spent entry
				}
			}

			return budget
		},
		func() *dto.BudgetList { return &dto.BudgetList{} },
	)
}

// MapBudgetLimitArrayToBudgetLimitList converts client.BudgetLimitArray to BudgetLimitList DTO.
// Maps budget limit data with amounts and currency information.
// Returns nil if the input is nil
func MapBudgetLimitArrayToBudgetLimitList(limitArray *client.BudgetLimitArray) *dto.BudgetLimitList {
	if limitArray == nil {
		return nil
	}

	limitList := &dto.BudgetLimitList{
		Data: make([]dto.BudgetLimit, 0),
	}

	// Map budget limit data
	for _, limitRead := range limitArray.Data {
		limit := dto.BudgetLimit{
			Id:             limitRead.Id,
			Amount:         limitRead.Attributes.Amount,
			Start:          limitRead.Attributes.Start,
			End:            limitRead.Attributes.End,
			BudgetId:       GetStringValue(limitRead.Attributes.BudgetId),
			CurrencyCode:   GetStringValue(limitRead.Attributes.CurrencyCode),
			CurrencySymbol: GetStringValue(limitRead.Attributes.CurrencySymbol),
			Spent:          []dto.BudgetSpent{},
		}

		// Map spent information if available
		if limitRead.Attributes.Spent != nil {
			// Spent is a single string amount, not an array
			budgetSpent := dto.BudgetSpent{
				Sum:            *limitRead.Attributes.Spent,
				CurrencyCode:   GetStringValue(limitRead.Attributes.CurrencyCode),
				CurrencySymbol: GetStringValue(limitRead.Attributes.CurrencySymbol),
			}
			limit.Spent = append(limit.Spent, budgetSpent)
		}

		limitList.Data = append(limitList.Data, limit)
	}

	// Map pagination
	if limitArray.Meta.Pagination != nil {
		pagination := limitArray.Meta.Pagination
		limitList.Pagination = dto.Pagination{
			Count:       GetIntValue(pagination.Count),
			Total:       GetIntValue(pagination.Total),
			CurrentPage: GetIntValue(pagination.CurrentPage),
			PerPage:     GetIntValue(pagination.PerPage),
			TotalPages:  GetIntValue(pagination.TotalPages),
		}
	}

	return limitList
}