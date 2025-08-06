package mappers

import "github.com/dezer32/mcp-firefly-iii/pkg/client"

// GetIntValue safely dereferences an int pointer
func GetIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// GetStringValue safely dereferences a string pointer
func GetStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// GetFloat32Value safely dereferences a float32 pointer
func GetFloat32Value(ptr *float32) float32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// GetBoolValue safely dereferences a bool pointer
func GetBoolValue(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

// GetAccountTypeString converts AccountTypeProperty pointer to string
func GetAccountTypeString(atp *client.AccountTypeProperty) string {
	if atp == nil {
		return ""
	}
	return string(*atp)
}