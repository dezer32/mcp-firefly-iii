package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapInsightGroupArrayToInsightCategoryList maps a client.InsightGroup to a dto.InsightCategoryResponse
func MapInsightGroupArrayToInsightCategoryList(insightGroup *client.InsightGroup) *dto.InsightCategoryResponse {
	responseBuilder := dto.NewInsightCategoryResponseBuilder()
	
	if insightGroup == nil {
		result := responseBuilder.Build()
		return &result
	}

	for _, group := range *insightGroup {
		// Each group has an ID, name, difference and currency code
		entryBuilder := dto.NewInsightCategoryEntryBuilder()
		
		if group.Id != nil {
			entryBuilder = entryBuilder.WithId(*group.Id)
		}
		if group.Name != nil {
			entryBuilder = entryBuilder.WithName(*group.Name)
		}
		if group.Difference != nil {
			entryBuilder = entryBuilder.WithAmount(*group.Difference)
		}
		if group.CurrencyCode != nil {
			entryBuilder = entryBuilder.WithCurrencyCode(*group.CurrencyCode)
		}
		
		responseBuilder = responseBuilder.AddEntry(entryBuilder.Build())
	}

	result := responseBuilder.Build()
	return &result
}

// MapInsightTotalArrayToInsightTotalList maps a client.InsightTotal to a dto.InsightTotalResponse
func MapInsightTotalArrayToInsightTotalList(insightTotal *client.InsightTotal) *dto.InsightTotalResponse {
	responseBuilder := dto.NewInsightTotalResponseBuilder()
	
	if insightTotal == nil {
		result := responseBuilder.Build()
		return &result
	}

	for _, totalEntry := range *insightTotal {
		// Each entry has difference and currency code
		entryBuilder := dto.NewInsightTotalEntryBuilder()
		
		if totalEntry.Difference != nil {
			entryBuilder = entryBuilder.WithAmount(*totalEntry.Difference)
		}
		if totalEntry.CurrencyCode != nil {
			entryBuilder = entryBuilder.WithCurrencyCode(*totalEntry.CurrencyCode)
		}
		
		responseBuilder = responseBuilder.AddEntry(entryBuilder.Build())
	}

	result := responseBuilder.Build()
	return &result
}