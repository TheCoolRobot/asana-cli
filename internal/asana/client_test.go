package asana

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "test-token-12345",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: false, // Error happens during API call, not initialization
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.token)
			if client == nil {
				t.Error("NewClient returned nil")
			}
			if client.baseURL != "https://app.asana.com/api/1.0" {
				t.Errorf("unexpected baseURL: %s", client.baseURL)
			}
		})
	}
}

func TestClientDoRequiresToken(t *testing.T) {
	client := NewClient("")
	
	_, err := client.do("GET", "/users/me", nil)
	if err == nil {
		t.Error("expected error when token is empty")
	}
}