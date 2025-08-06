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
			return dto.Category{
				Id:    categoryRead.Id,
				Name:  categoryRead.Attributes.Name,
				Notes: categoryRead.Attributes.Notes,
			}
		},
		func() *dto.CategoryList { return &dto.CategoryList{} },
	)
}