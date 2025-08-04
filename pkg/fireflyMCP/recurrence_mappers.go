package fireflyMCP

import (
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
)

// mapRecurrenceRepetitionToDTO converts client RecurrenceRepetition to DTO
func mapRecurrenceRepetitionToDTO(rep *client.RecurrenceRepetition) RecurrenceRepetition {
	if rep == nil {
		return RecurrenceRepetition{}
	}

	result := RecurrenceRepetition{
		Id:          getStringValue(rep.Id),
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

// mapRecurrenceTransactionToDTO converts client RecurrenceTransaction to DTO
func mapRecurrenceTransactionToDTO(trans *client.RecurrenceTransaction) RecurrenceTransaction {
	if trans == nil {
		return RecurrenceTransaction{}
	}

	return RecurrenceTransaction{
		Id:              getStringValue(trans.Id),
		Description:     trans.Description,
		Amount:          trans.Amount,
		CurrencyCode:    getStringValue(trans.CurrencyCode),
		CategoryId:      trans.CategoryId,
		CategoryName:    trans.CategoryName,
		BudgetId:        trans.BudgetId,
		BudgetName:      trans.BudgetName,
		SourceId:        getStringValue(trans.SourceId),
		SourceName:      getStringValue(trans.SourceName),
		DestinationId:   getStringValue(trans.DestinationId),
		DestinationName: getStringValue(trans.DestinationName),
	}
}

// mapRecurrenceToRecurrence converts client.RecurrenceRead to Recurrence DTO
func mapRecurrenceToRecurrence(recurrenceRead *client.RecurrenceRead) *Recurrence {
	if recurrenceRead == nil {
		return nil
	}

	recurrence := &Recurrence{
		Id:              recurrenceRead.Id,
		Type:            string(getRecurrenceType(recurrenceRead.Attributes.Type)),
		Title:           getStringValue(recurrenceRead.Attributes.Title),
		Description:     getStringValue(recurrenceRead.Attributes.Description),
		FirstDate:       time.Time{},
		LatestDate:      nil,
		RepeatUntil:     nil,
		NrOfRepetitions: nil,
		ApplyRules:      false,
		Active:          false,
		Notes:           recurrenceRead.Attributes.Notes,
		Repetitions:     []RecurrenceRepetition{},
		Transactions:    []RecurrenceTransaction{},
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
			recurrence.Repetitions = append(recurrence.Repetitions, mapRecurrenceRepetitionToDTO(&rep))
		}
	}

	// Map transactions
	if recurrenceRead.Attributes.Transactions != nil {
		for _, trans := range *recurrenceRead.Attributes.Transactions {
			recurrence.Transactions = append(recurrence.Transactions, mapRecurrenceTransactionToDTO(&trans))
		}
	}

	return recurrence
}

// mapRecurrenceArrayToRecurrenceList converts client.RecurrenceArray to RecurrenceList DTO
func mapRecurrenceArrayToRecurrenceList(recurrenceArray *client.RecurrenceArray) RecurrenceList {
	recurrenceList := RecurrenceList{
		Data: []Recurrence{},
	}

	if recurrenceArray == nil {
		return recurrenceList
	}

	// Reinitialize with correct capacity
	recurrenceList.Data = make([]Recurrence, 0, len(recurrenceArray.Data))

	// Map recurrence data
	for _, recurrenceRead := range recurrenceArray.Data {
		if mappedRecurrence := mapRecurrenceToRecurrence(&recurrenceRead); mappedRecurrence != nil {
			recurrenceList.Data = append(recurrenceList.Data, *mappedRecurrence)
		}
	}

	// Map pagination
	if recurrenceArray.Meta.Pagination != nil {
		pagination := recurrenceArray.Meta.Pagination
		recurrenceList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return recurrenceList
}
