package mappers

import (
	"reflect"

	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// ItemMapper is a function type that maps from source type S to destination type D
type ItemMapper[S any, D any] func(S) D

// MapArrayToList is a generic function that maps array data to list DTOs with pagination.
// It significantly reduces code duplication across all mapper functions.
//
// This function uses reflection to handle pagination mapping, which allows it to work
// with any pagination structure that has the standard fields.
//
// Parameters:
//   - arrayData: the source array struct (e.g., *client.AccountArray)
//   - mapItem: function to map individual items from source to destination type
//   - createList: function to create a new list instance
//
// Returns:
//   - L: mapped list with data and pagination
//
// Example usage:
//
//	func MapAccountArrayToAccountList(accountArray *client.AccountArray) *dto.AccountList {
//	    if accountArray == nil {
//	        return nil
//	    }
//	    return MapArrayToList(
//	        accountArray,
//	        func(r client.AccountRead) dto.Account {
//	            return dto.Account{
//	                Id:     r.Id,
//	                Active: r.Attributes.Active != nil && *r.Attributes.Active,
//	                Name:   r.Attributes.Name,
//	                Notes:  r.Attributes.Notes,
//	                Type:   string(r.Attributes.Type),
//	            }
//	        },
//	        func() *dto.AccountList { return &dto.AccountList{} },
//	    )
//	}
func MapArrayToList[A any, S any, D any, L any](
	arrayData A,
	mapItem ItemMapper[S, D],
	createList func() L,
) L {
	list := createList()

	// Use reflection to access the Data field
	arrayValue := reflect.ValueOf(arrayData)
	if !arrayValue.IsValid() || arrayValue.IsNil() {
		return list
	}

	// If it's a pointer, get the underlying value
	if arrayValue.Kind() == reflect.Ptr {
		arrayValue = arrayValue.Elem()
	}

	// Get the Data field
	dataField := arrayValue.FieldByName("Data")
	if dataField.IsValid() && dataField.Kind() == reflect.Slice {
		// Convert to slice and map items
		dataSlice := dataField.Interface().([]S)
		mappedData := make([]D, 0, len(dataSlice))
		for _, item := range dataSlice {
			mappedData = append(mappedData, mapItem(item))
		}

		// Set the Data field in the list
		listValue := reflect.ValueOf(list)
		if listValue.Kind() == reflect.Ptr {
			listValue = listValue.Elem()
		}
		listDataField := listValue.FieldByName("Data")
		if listDataField.IsValid() && listDataField.CanSet() {
			listDataField.Set(reflect.ValueOf(mappedData))
		}
	}

	// Get and map pagination
	metaField := arrayValue.FieldByName("Meta")
	if metaField.IsValid() {
		paginationField := metaField.FieldByName("Pagination")
		if paginationField.IsValid() && !paginationField.IsNil() {
			pagination := mapPaginationUsingReflection(paginationField.Interface())

			// Set pagination in the list
			listValue := reflect.ValueOf(list)
			if listValue.Kind() == reflect.Ptr {
				listValue = listValue.Elem()
			}
			listPaginationField := listValue.FieldByName("Pagination")
			if listPaginationField.IsValid() && listPaginationField.CanSet() {
				listPaginationField.Set(reflect.ValueOf(pagination))
			}
		}
	}

	return list
}

// mapPaginationUsingReflection maps pagination using reflection to handle any pagination struct
func mapPaginationUsingReflection(pagination interface{}) dto.Pagination {
	if pagination == nil {
		return dto.NewPaginationBuilder().Build()
	}

	paginationValue := reflect.ValueOf(pagination)

	// Handle pointer
	if paginationValue.Kind() == reflect.Ptr {
		if paginationValue.IsNil() {
			return dto.NewPaginationBuilder().Build()
		}
		paginationValue = paginationValue.Elem()
	}

	// Map fields using builder
	builder := dto.NewPaginationBuilder()
	
	if field := paginationValue.FieldByName("Count"); field.IsValid() {
		if ptr, ok := field.Interface().(*int); ok {
			builder = builder.WithCount(GetIntValue(ptr))
		}
	}
	if field := paginationValue.FieldByName("Total"); field.IsValid() {
		if ptr, ok := field.Interface().(*int); ok {
			builder = builder.WithTotal(GetIntValue(ptr))
		}
	}
	if field := paginationValue.FieldByName("CurrentPage"); field.IsValid() {
		if ptr, ok := field.Interface().(*int); ok {
			builder = builder.WithCurrentPage(GetIntValue(ptr))
		}
	}
	if field := paginationValue.FieldByName("PerPage"); field.IsValid() {
		if ptr, ok := field.Interface().(*int); ok {
			builder = builder.WithPerPage(GetIntValue(ptr))
		}
	}
	if field := paginationValue.FieldByName("TotalPages"); field.IsValid() {
		if ptr, ok := field.Interface().(*int); ok {
			builder = builder.WithTotalPages(GetIntValue(ptr))
		}
	}

	return builder.Build()
}

// MapPaginationToDTO is a simple helper that maps pagination fields
// This is exported for use in mappers that need custom pagination handling
func MapPaginationToDTO(pagination interface{}) dto.Pagination {
	return mapPaginationUsingReflection(pagination)
}