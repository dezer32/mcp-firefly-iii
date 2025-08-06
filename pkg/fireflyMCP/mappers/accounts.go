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
	return dto.NewAccountBuilder().
	WithId(accountRead.Id).
	WithActive(accountRead.Attributes.Active != nil && *accountRead.Attributes.Active).
	WithName(accountRead.Attributes.Name).
	WithNotes(accountRead.Attributes.Notes).
	WithType(string(accountRead.Attributes.Type)).
	 Build()
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

	account := dto.NewAccountBuilder().
		WithId(accountSingle.Data.Id).
		WithActive(accountSingle.Data.Attributes.Active != nil && *accountSingle.Data.Attributes.Active).
		WithName(accountSingle.Data.Attributes.Name).
		WithNotes(accountSingle.Data.Attributes.Notes).
		WithType(string(accountSingle.Data.Attributes.Type)).
		Build()
	return &account
}

// MapAccountReadToAccount converts client.AccountRead to Account DTO.
// It maps individual account attributes to the simplified DTO format.
// Returns nil if the input is nil
func MapAccountReadToAccount(accountRead *client.AccountRead) *dto.Account {
	if accountRead == nil {
		return nil
	}

	account := dto.NewAccountBuilder().
		WithId(accountRead.Id).
		WithActive(accountRead.Attributes.Active != nil && *accountRead.Attributes.Active).
		WithName(accountRead.Attributes.Name).
		WithNotes(accountRead.Attributes.Notes).
		WithType(string(accountRead.Attributes.Type)).
		Build()
	return &account
}