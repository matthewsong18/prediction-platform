package main

import (
	"os"
	"testing"
)

func TestLoadConfig_EncryptionKey(t *testing.T) {
	// Helper to set env vars
	setEnv := func(env map[string]string) {
		os.Clearenv()
		for k, v := range env {
			os.Setenv(k, v)
		}
	}

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "Valid Key",
			env: map[string]string{
				"GUILD_ID":       "123",
				"TOKEN":          "abc",
				"APP_ID":         "456",
				"DB_PATH":        "test.db",
				"ENCRYPTION_KEY": "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f", // 32 bytes hex
			},
			wantErr: false,
		},
		{
			name: "Missing Key",
			env: map[string]string{
				"GUILD_ID": "123",
				"TOKEN":    "abc",
				"APP_ID":   "456",
				"DB_PATH":  "test.db",
			},
			wantErr: true,
		},
		{
			name: "Invalid Hex",
			env: map[string]string{
				"GUILD_ID":       "123",
				"TOKEN":          "abc",
				"APP_ID":         "456",
				"DB_PATH":        "test.db",
				"ENCRYPTION_KEY": "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", // Invalid hex chars
			},
			wantErr: true,
		},
		{
			name: "Wrong Length (Too Short)",
			env: map[string]string{
				"GUILD_ID":       "123",
				"TOKEN":          "abc",
				"APP_ID":         "456",
				"DB_PATH":        "test.db",
				"ENCRYPTION_KEY": "deadbeef",
			},
			wantErr: true,
		},
		{
			name: "Wrong Length (Too Long)",
			env: map[string]string{
				"GUILD_ID":       "123",
				"TOKEN":          "abc",
				"APP_ID":         "456",
				"DB_PATH":        "test.db",
				"ENCRYPTION_KEY": "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f00",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.env)
			_, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
