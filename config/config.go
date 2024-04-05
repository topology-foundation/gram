package config

import (
	"encoding/json"
	"flag"
	"os"
)

var (
	configPath = flag.String("config", "config.json", "The path to the configuration file")
)

type ExecutionConfig struct {
}

type NetworkConfig struct {
	Namespace string   `json:"namespace"`
	Topics    []string `json:"topics"`
	MaxPeers  int      `json:"maxPeers"`
	Port      int      `json:"port"`
}

type StorageConfig struct {
	DatabasePath string `json:"databasePath"`
}

type GrpcConfig struct {
	Port int `json:"port"`
}

type AppConfig struct {
	Execution ExecutionConfig `json:"executionConfig"`
	Network   NetworkConfig   `json:"networkConfig"`
	Storage   StorageConfig   `json:"storageConfig"`
	Grpc      GrpcConfig      `json:"grpcConfig"`
}

func LoadConfig() (*AppConfig, error) {
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &AppConfig{}
	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
