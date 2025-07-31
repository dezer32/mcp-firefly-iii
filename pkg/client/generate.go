//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate client,models,embedded-spec -package client -o client.go ../../resources/firefly-iii-6.2.21-v1.yaml

package client

// This file contains the go:generate directive for regenerating the Firefly III client code.
// To regenerate the client, run: go generate ./pkg/client