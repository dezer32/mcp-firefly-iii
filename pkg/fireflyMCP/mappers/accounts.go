package mappers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// MapAccountArrayToAccountList converts client.AccountArray to AccountList DTO.
// It maps account data and pagination information from the API response.
// Returns nil if the input is nil
func MapAccountArrayToAccountList(accountArray *client.AccountArray) *dto.AccountList {
	if accountArray == nil {
		return nil
	}

	return MapArrayToList(
		accountArray,
		func(accountRead client.AccountRead) dto.Account {
			return dto.Account{
				Id:     accountRead.Id,
				Active: accountRead.Attributes.Active != nil && *accountRead.Attributes.Active,
				Name:   accountRead.Attributes.Name,
				Notes:  accountRead.Attributes.Notes,
				Type:   string(accountRead.Attributes.Type),
			}
		},
		func() *dto.AccountList { return &dto.AccountList{} },
	)
}

// MapAccountSingleToAccount converts client.AccountSingle to Account DTO.
// It extracts account data from the single account API response.
// Returns nil if the input is nil
func MapAccountSingleToAccount(accountSingle *client.AccountSingle) *dto.Account {
	if accountSingle == nil {
		return nil
	}

	return &dto.Account{
		Id:     accountSingle.Data.Id,
		Active: accountSingle.Data.Attributes.Active != nil && *accountSingle.Data.Attributes.Active,
		Name:   accountSingle.Data.Attributes.Name,
		Notes:  accountSingle.Data.Attributes.Notes,
		Type:   string(accountSingle.Data.Attributes.Type),
	}
}

// MapAccountReadToAccount converts client.AccountRead to Account DTO.
// It maps individual account attributes to the simplified DTO format.
// Returns nil if the input is nil
func MapAccountReadToAccount(accountRead *client.AccountRead) *dto.Account {
	if accountRead == nil {
		return nil
	}

	return &dto.Account{
		Id:     accountRead.Id,
		Active: accountRead.Attributes.Active != nil && *accountRead.Attributes.Active,
		Name:   accountRead.Attributes.Name,
		Notes:  accountRead.Attributes.Notes,
		Type:   string(accountRead.Attributes.Type),
	}
}