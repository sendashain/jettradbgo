package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port           string
	DataDir        string
	ClusterEnabled bool
	ClusterPort    string
	ReplicationFactor int
	ConsistencyLevel  string
}

// LoadConfig loads configuration from environment variables or uses defaults
func LoadConfig() *Config {
	return &Config{
		Port:              getEnvOrDefault("DB_PORT", "8080"),
		DataDir:           getEnvOrDefault("DB_DATA_DIR", "./data"),
		ClusterEnabled:    getEnvOrDefaultBool("CLUSTER_ENABLED", false),
		ClusterPort:       getEnvOrDefault("CLUSTER_PORT", "9090"),
		ReplicationFactor: getEnvOrDefaultInt("REPLICATION_FACTOR", 1),
		ConsistencyLevel:  getEnvOrDefault("CONSISTENCY_LEVEL", "quorum"),
	}
}

// Helper functions to get environment variables with defaults
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" || value == "1" {
			return true
		}
		return false
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}