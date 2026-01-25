package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

type Config struct {
	GuildID       string
	Token         string
	AppID         string
	DBPath        string
	EncryptionKey string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		GuildID:       os.Getenv("GUILD_ID"),
		Token:         os.Getenv("TOKEN"),
		AppID:         os.Getenv("APP_ID"),
		DBPath:        os.Getenv("DB_PATH"),
		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),
	}

	if cfg.GuildID == "" {
		return nil, fmt.Errorf("GUILD_ID environment variable is not set")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("TOKEN environment variable is not set")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("APP_ID environment variable is not set")
	}
	if cfg.DBPath == "" {
		return nil, fmt.Errorf("DB_PATH environment variable is not set")
	}

	if cfg.EncryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable is not set")
	}

	// Validate Encryption Key
	keyBytes, err := hex.DecodeString(cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be a valid hex string: %w", err)
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes (64 hex characters), got %d bytes", len(keyBytes))
	}

	return cfg, nil
}
