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

	recurrence := &dto.Recurrence{
		Id:              recurrenceRead.Id,
		Type:            string(GetRecurrenceType(recurrenceRead.Attributes.Type)),
		Title:           GetStringValue(recurrenceRead.Attributes.Title),
		Description:     GetStringValue(recurrenceRead.Attributes.Description),
		FirstDate:       time.Time{},
		LatestDate:      nil,
		RepeatUntil:     nil,
		NrOfRepetitions: nil,
		ApplyRules:      false,
		Active:          false,
		Notes:           recurrenceRead.Attributes.Notes,
		Repetitions:     []dto.RecurrenceRepetition{},
		Transactions:    []dto.RecurrenceTransaction{},
	}

	// Handle dates
	if recurrenceRead.Attributes.FirstDate != nil {
		recurrence.FirstDate = recurrenceRead.Attributes.FirstDate.Time
	}
	if recurrenceRead.Attributes.LatestDate != nil {
		recurrence.LatestDate = &recurrenceRead.Attributes.LatestDate.Time
	}
	if recurrenceRead.Attributes.RepeatUntil != nil {
		recurrence.RepeatUntil = &recurrenceRead.Attributes.RepeatUntil.Time
	}

	// Handle other fields
	if recurrenceRead.Attributes.NrOfRepetitions != nil {
		nr := int(*recurrenceRead.Attributes.NrOfRepetitions)
		recurrence.NrOfRepetitions = &nr
	}
	if recurrenceRead.Attributes.ApplyRules != nil {
		recurrence.ApplyRules = *recurrenceRead.Attributes.ApplyRules
	}
	if recurrenceRead.Attributes.Active != nil {
		recurrence.Active = *recurrenceRead.Attributes.Active
	}

	// Map repetitions
	if recurrenceRead.Attributes.Repetitions != nil {
		for _, rep := range *recurrenceRead.Attributes.Repetitions {
			recurrence.Repetitions = append(recurrence.Repetitions, MapRecurrenceRepetitionToDTO(&rep))
		}
	}

	// Map transactions
	if recurrenceRead.Attributes.Transactions != nil {
		for _, trans := range *recurrenceRead.Attributes.Transactions {
			recurrence.Transactions = append(recurrence.Transactions, MapRecurrenceTransactionToDTO(&trans))
		}
	}

	return recurrence
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
		recurrenceList.Pagination = dto.Pagination{
			Count:       GetIntValue(pagination.Count),
			Total:       GetIntValue(pagination.Total),
			CurrentPage: GetIntValue(pagination.CurrentPage),
			PerPage:     GetIntValue(pagination.PerPage),
			TotalPages:  GetIntValue(pagination.TotalPages),
		}
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