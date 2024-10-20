package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port              int    `json:"port"`
	ChunkSize         int64  `json:"chunk_size"`
	ReplicationFactor int    `json:"replication_factor"`
	StoragePath       string `json:"storage_path"`
	EncryptionKeyPath string `json:"encryption_key_path"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetEncryptionKey() ([]byte, error) {
	key, err := os.ReadFile(c.EncryptionKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encryption key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes long")
	}
	return key, nil
}