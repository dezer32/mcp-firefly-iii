package mappers

import (
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapRecurrenceRepetitionToDTO converts client RecurrenceRepetition to DTO
func MapRecurrenceRepetitionToDTO(rep *client.RecurrenceRepetition) dto.RecurrenceRepetition {
	if rep == nil {
		return dto.RecurrenceRepetition{}
	}

	result := dto.RecurrenceRepetition{
		Id:          GetStringValue(rep.Id),
		Type:        string(rep.Type),
		Moment:      rep.Moment,
		Skip:        0,
		Weekend:     0,
		Description: rep.Description,
	}

	if rep.Skip != nil {
		result.Skip = int(*rep.Skip)
	}
	if rep.Weekend != nil {
		result.Weekend = int(*rep.Weekend)
	}

	return result
}

// MapRecurrenceTransactionToDTO converts client RecurrenceTransaction to DTO
func MapRecurrenceTransactionToDTO(trans *client.RecurrenceTransaction) dto.RecurrenceTransaction {
	if trans == nil {
		return dto.RecurrenceTransaction{}
	}

	return dto.RecurrenceTransaction{
		Id:              GetStringValue(trans.Id),
		Description:     trans.Description,
		Amount:          trans.Amount,
		CurrencyCode:    GetStringValue(trans.CurrencyCode),
		CategoryId:      trans.CategoryId,
		CategoryName:    trans.CategoryName,
		BudgetId:        trans.BudgetId,
		BudgetName:      trans.BudgetName,
		SourceId:        GetStringValue(trans.SourceId),
		SourceName:      GetStringValue(trans.SourceName),
		DestinationId:   GetStringValue(trans.DestinationId),
		DestinationName: GetStringValue(trans.DestinationName),
	}
}

// MapRecurrenceToRecurrence converts client.RecurrenceRead to Recurrence DTO
func MapRecurrenceToRecurrence(recurrenceRead *client.RecurrenceRead) *dto.Recurrence {
	if recurrenceRead == nil {
		return nil
	}

	builder := dto.NewRecurrenceBuilder().
		WithId(recurrenceRead.Id).
		WithType(string(GetRecurrenceType(recurrenceRead.Attributes.Type))).
		WithTitle(GetStringValue(recurrenceRead.Attributes.Title)).
		WithDescription(GetStringValue(recurrenceRead.Attributes.Description)).
		WithNotes(recurrenceRead.Attributes.Notes)

	// Handle dates
	firstDate := time.Time{}
	if recurrenceRead.Attributes.FirstDate != nil {
		firstDate = recurrenceRead.Attributes.FirstDate.Time
	}
	builder = builder.WithFirstDate(firstDate)
	
	if recurrenceRead.Attributes.LatestDate != nil {
		builder = builder.WithLatestDate(&recurrenceRead.Attributes.LatestDate.Time)
	}
	if recurrenceRead.Attributes.RepeatUntil != nil {
		builder = builder.WithRepeatUntil(&recurrenceRead.Attributes.RepeatUntil.Time)
	}

	// Handle other fields
	if recurrenceRead.Attributes.NrOfRepetitions != nil {
		nr := int(*recurrenceRead.Attributes.NrOfRepetitions)
		builder = builder.WithNrOfRepetitions(&nr)
	}
	
	applyRules := false
	if recurrenceRead.Attributes.ApplyRules != nil {
		applyRules = *recurrenceRead.Attributes.ApplyRules
	}
	builder = builder.WithApplyRules(applyRules)
	
	active := false
	if recurrenceRead.Attributes.Active != nil {
		active = *recurrenceRead.Attributes.Active
	}
	builder = builder.WithActive(active)

	// Map repetitions
	repetitions := []dto.RecurrenceRepetition{}
	if recurrenceRead.Attributes.Repetitions != nil {
		for _, rep := range *recurrenceRead.Attributes.Repetitions {
			repetitions = append(repetitions, MapRecurrenceRepetitionToDTO(&rep))
		}
	}
	builder = builder.WithRepetitions(repetitions)

	// Map transactions
	transactions := []dto.RecurrenceTransaction{}
	if recurrenceRead.Attributes.Transactions != nil {
		for _, trans := range *recurrenceRead.Attributes.Transactions {
			transactions = append(transactions, MapRecurrenceTransactionToDTO(&trans))
		}
	}
	builder = builder.WithTransactions(transactions)

	recurrence := builder.Build()
	return &recurrence
}

// MapRecurrenceArrayToRecurrenceList converts client.RecurrenceArray to RecurrenceList DTO
func MapRecurrenceArrayToRecurrenceList(recurrenceArray *client.RecurrenceArray) dto.RecurrenceList {
	recurrenceList := dto.RecurrenceList{
		Data: []dto.Recurrence{},
	}

	if recurrenceArray == nil {
		return recurrenceList
	}

	// Reinitialize with correct capacity
	recurrenceList.Data = make([]dto.Recurrence, 0, len(recurrenceArray.Data))

	// Map recurrence data
	for _, recurrenceRead := range recurrenceArray.Data {
		if mappedRecurrence := MapRecurrenceToRecurrence(&recurrenceRead); mappedRecurrence != nil {
			recurrenceList.Data = append(recurrenceList.Data, *mappedRecurrence)
		}
	}

	// Map pagination
	if recurrenceArray.Meta.Pagination != nil {
		pagination := recurrenceArray.Meta.Pagination
		recurrenceList.Pagination = dto.NewPaginationBuilder().
			WithCount(GetIntValue(pagination.Count)).
			WithTotal(GetIntValue(pagination.Total)).
			WithCurrentPage(GetIntValue(pagination.CurrentPage)).
			WithPerPage(GetIntValue(pagination.PerPage)).
			WithTotalPages(GetIntValue(pagination.TotalPages)).
			Build()
	}

	return recurrenceList
}


// GetRecurrenceType helper function
func GetRecurrenceType(t *client.RecurrenceTransactionType) client.RecurrenceTransactionType {
	if t == nil {
		return ""
	}
	return *t
}