package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapCategoryArrayToCategoryList converts client.CategoryArray to CategoryList DTO
func MapCategoryArrayToCategoryList(categoryArray *client.CategoryArray) *dto.CategoryList {
	if categoryArray == nil {
		return nil
	}

	return MapArrayToList(
		categoryArray,
		func(categoryRead client.CategoryRead) dto.Category {
			return dto.NewCategoryBuilder().
				WithId(categoryRead.Id).
				WithName(categoryRead.Attributes.Name).
				WithNotes(categoryRead.Attributes.Notes).
				Build()
		},
		func() *dto.CategoryList { return &dto.CategoryList{} },
	)
}