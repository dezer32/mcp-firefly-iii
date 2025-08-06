package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapTagArrayToTagList maps a client.TagArray to a dto.TagList
func MapTagArrayToTagList(tagArray *client.TagArray) *dto.TagList {
	if tagArray == nil {
		return &dto.TagList{
			Data:       []dto.Tag{},
			Pagination: dto.NewPaginationBuilder().Build(),
		}
	}

	return MapArrayToList(
		tagArray,
		func(tagRead client.TagRead) dto.Tag {
			return dto.NewTagBuilder().
				WithId(tagRead.Id).
				WithTag(tagRead.Attributes.Tag).
				WithDescription(tagRead.Attributes.Description).
				Build()
		},
		func() *dto.TagList { return &dto.TagList{} },
	)
}

// MapTagReadToTag maps a client.TagRead to a dto.Tag
func MapTagReadToTag(tagRead *client.TagRead) *dto.Tag {
	if tagRead == nil {
		emptyTag := dto.NewTagBuilder().Build()
		return &emptyTag
	}

	tag := dto.NewTagBuilder().
		WithId(tagRead.Id).
		WithTag(tagRead.Attributes.Tag).
		WithDescription(tagRead.Attributes.Description).
		Build()

	return &tag
}