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
			Pagination: dto.Pagination{},
		}
	}

	return MapArrayToList(
		tagArray,
		func(tagRead client.TagRead) dto.Tag {
			return dto.Tag{
				Id:          tagRead.Id,
				Tag:         tagRead.Attributes.Tag,
				Description: tagRead.Attributes.Description,
			}
		},
		func() *dto.TagList { return &dto.TagList{} },
	)
}

// MapTagReadToTag maps a client.TagRead to a dto.Tag
func MapTagReadToTag(tagRead *client.TagRead) *dto.Tag {
	if tagRead == nil {
		return &dto.Tag{}
	}

	tag := &dto.Tag{
		Id:          tagRead.Id,
		Tag:         tagRead.Attributes.Tag,
		Description: tagRead.Attributes.Description,
	}

	return tag
}