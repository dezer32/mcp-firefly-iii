package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapInsightGroupArrayToInsightCategoryList maps a client.InsightGroup to a dto.InsightCategoryResponse
func MapInsightGroupArrayToInsightCategoryList(insightGroup *client.InsightGroup) *dto.InsightCategoryResponse {
	if insightGroup == nil {
		return &dto.InsightCategoryResponse{
			Entries: []dto.InsightCategoryEntry{},
		}
	}

	entries := make([]dto.InsightCategoryEntry, 0)
	for _, group := range *insightGroup {
		// Each group has an ID, name, difference and currency code
		entry := dto.InsightCategoryEntry{}
		
		if group.Id != nil {
			entry.Id = *group.Id
		}
		if group.Name != nil {
			entry.Name = *group.Name
		}
		if group.Difference != nil {
			entry.Amount = *group.Difference
		}
		if group.CurrencyCode != nil {
			entry.CurrencyCode = *group.CurrencyCode
		}
		
		entries = append(entries, entry)
	}

	return &dto.InsightCategoryResponse{
		Entries: entries,
	}
}

// MapInsightTotalArrayToInsightTotalList maps a client.InsightTotal to a dto.InsightTotalResponse
func MapInsightTotalArrayToInsightTotalList(insightTotal *client.InsightTotal) *dto.InsightTotalResponse {
	if insightTotal == nil {
		return &dto.InsightTotalResponse{
			Entries: []dto.InsightTotalEntry{},
		}
	}

	entries := make([]dto.InsightTotalEntry, 0)
	for _, totalEntry := range *insightTotal {
		// Each entry has difference and currency code
		entry := dto.InsightTotalEntry{}
		
		if totalEntry.Difference != nil {
			entry.Amount = *totalEntry.Difference
		}
		if totalEntry.CurrencyCode != nil {
			entry.CurrencyCode = *totalEntry.CurrencyCode
		}
		
		entries = append(entries, entry)
	}

	return &dto.InsightTotalResponse{
		Entries: entries,
	}
}