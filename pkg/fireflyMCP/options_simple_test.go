package fireflyMCP

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionalOptions_Simple(t *testing.T) {
	t.Run("ServerWithMinimalConfig", func(t *testing.T) {
		server, err := NewServerWithOptions(
			WithBaseURL("https://test.example.com"),
			WithAPIToken("test-token-123"),
		)
		
		require.NoError(t, err)
		require.NotNil(t, server)
		assert.Equal(t, "https://test.example.com", server.Config.Server.URL)
		assert.Equal(t, "test-token-123", server.Config.API.Token)
	})
	
	t.Run("ServerWithTimeout", func(t *testing.T) {
		server, err := NewServerWithOptions(
			WithBaseURL("https://test.example.com"),
			WithAPIToken("test-token-123"),
			WithTimeout(45*time.Second),
		)
		
		require.NoError(t, err)
		require.NotNil(t, server)
		assert.Equal(t, 45, server.Config.Client.Timeout)
	})
	
	t.Run("ServerWithMCPInfo", func(t *testing.T) {
		server, err := NewServerWithOptions(
			WithBaseURL("https://test.example.com"),
			WithAPIToken("test-token-123"),
			WithMCPInfo("custom-mcp", "2.0.0"),
		)
		
		require.NoError(t, err)
		require.NotNil(t, server)
		assert.Equal(t, "custom-mcp", server.Config.MCP.Name)
		assert.Equal(t, "2.0.0", server.Config.MCP.Version)
	})
	
	t.Run("BackwardCompatibilityWithNewServer", func(t *testing.T) {
		config := &Config{}
		config.Server.URL = "https://test.example.com"
		config.API.Token = "test-token-123"
		config.Client.Timeout = 30
		config.MCP.Name = "test-mcp"
		config.MCP.Version = "1.0.0"
		
		server, err := NewServer(config)
		
		require.NoError(t, err)
		require.NotNil(t, server)
		assert.Equal(t, config.Server.URL, server.Config.Server.URL)
		assert.Equal(t, config.API.Token, server.Config.API.Token)
	})
	
	t.Run("ClientWithMinimalConfig", func(t *testing.T) {
		client, err := NewFireflyClient(
			WithClientBaseURL("https://test.example.com"),
			WithClientAPIToken("test-token-123"),
		)
		
		require.NoError(t, err)
		require.NotNil(t, client)
	})
	
	t.Run("ErrorOnMissingRequiredFields", func(t *testing.T) {
		// Missing API token
		server, err := NewServerWithOptions(
			WithBaseURL("https://test.example.com"),
		)
		
		assert.Error(t, err)
		assert.Nil(t, server)
		
		// Missing base URL
		server, err = NewServerWithOptions(
			WithAPIToken("test-token-123"),
		)
		
		assert.Error(t, err)
		assert.Nil(t, server)
	})
}