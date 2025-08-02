package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dezer32/firefly-iii/pkg/client"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		URL string `yaml:"url"`
	} `yaml:"server"`
	API struct {
		Token string `yaml:"token"`
	} `yaml:"api"`
	Client struct {
		Timeout int `yaml:"timeout"`
	} `yaml:"client"`
	Limits struct {
		Accounts     int `yaml:"accounts"`
		Transactions int `yaml:"transactions"`
	} `yaml:"limits"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// AuthTransport adds authentication to HTTP requests
type AuthTransport struct {
	Token string
	Base  http.RoundTripper
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

func main() {
	// Load configuration
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create HTTP client with authentication and timeout
	httpClient := &http.Client{
		Transport: &AuthTransport{
			Token: config.API.Token,
		},
		Timeout: time.Duration(config.Client.Timeout) * time.Second,
	}

	// Create Firefly III client
	fireflyClient, err := client.NewClientWithResponses(
		config.Server.URL,
		client.WithHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Get system information
	fmt.Println("=== Getting System Information ===")
	aboutResp, err := fireflyClient.GetAboutWithResponse(ctx, &client.GetAboutParams{})
	if err != nil {
		log.Printf("Error getting about info: %v", err)
	} else {
		fmt.Printf("Response status: %s\n", aboutResp.Status())
		if aboutResp.StatusCode() == 200 {
			fmt.Println("Successfully connected to Firefly III!")
		}
	}

	// Example 2: List accounts
	fmt.Println("\n=== Listing Accounts ===")
	accountsResp, err := fireflyClient.ListAccountWithResponse(ctx, &client.ListAccountParams{
		Limit: &[]int32{int32(config.Limits.Accounts)}[0], // Get accounts from config
	})
	if err != nil {
		log.Printf("Error listing accounts: %v", err)
	} else {
		fmt.Printf("Accounts response status: %s\n", accountsResp.Status())
		if accountsResp.StatusCode() == 200 {
			fmt.Println("Successfully retrieved accounts list!")
		}
	}

	// Example 3: List transactions
	fmt.Println("\n=== Listing Recent Transactions ===")
	transactionsResp, err := fireflyClient.ListTransactionWithResponse(ctx, &client.ListTransactionParams{
		Limit: &[]int32{int32(config.Limits.Transactions)}[0], // Get transactions from config
	})
	if err != nil {
		log.Printf("Error listing transactions: %v", err)
	} else {
		fmt.Printf("Transactions response status: %s\n", transactionsResp.Status())
		if transactionsResp.StatusCode() == 200 {
			fmt.Println("Successfully retrieved transactions list!")
		}
	}

	// Example 4: Search accounts
	fmt.Println("\n=== Searching for Accounts ===")
	searchResp, err := fireflyClient.SearchAccountsWithResponse(ctx, &client.SearchAccountsParams{
		Query: "checking",
		Field: client.AccountSearchFieldFilterName,
		Limit: &[]int32{5}[0],
	})
	if err != nil {
		log.Printf("Error searching accounts: %v", err)
	} else {
		fmt.Printf("Search response status: %s\n", searchResp.Status())
		if searchResp.StatusCode() == 200 {
			fmt.Println("Successfully searched accounts!")
		}
	}

	// Example 5: Get basic summary
	fmt.Println("\n=== Getting Basic Summary ===")
	summaryResp, err := fireflyClient.GetBasicSummaryWithResponse(ctx, &client.GetBasicSummaryParams{})
	if err != nil {
		log.Printf("Error getting summary: %v", err)
	} else {
		fmt.Printf("Summary response status: %s\n", summaryResp.Status())
		if summaryResp.StatusCode() == 200 {
			fmt.Println("Successfully retrieved basic summary!")
		}
	}

	fmt.Println("\n=== Example completed ===")
	fmt.Println("To use this example:")
	fmt.Println("1. Edit config.yaml file:")
	fmt.Println("   - Replace 'your-firefly-iii-instance.com' with your actual Firefly III URL")
	fmt.Println("   - Replace 'your-api-token-here' with your actual API token")
	fmt.Println("   - Adjust timeout and limits as needed")
	fmt.Println("2. Run: go run example/main.go")
}
