package config

import (
	"fmt"
	"os"
)

type ConnectionConfig struct {
	Driver string
	Host   string
	Port   string
	User   string
	Pass   string
	Name   string
}

func GetConfig() ConnectionConfig {
	return ConnectionConfig{
		Driver: getEnv("DB_DRIVER", "postgres"),
		Host:   getEnv("DB_HOST", "localhost"),
		Port:   getEnv("DB_PORT", "54329"),
		User:   getEnv("DB_USERNAME", "postgres"),
		Pass:   getEnv("DB_PASSWORD", "postgres"),
		Name:   getEnv("DB_DATABASE", "test"),
	}
}

func (cc *ConnectionConfig) GetDSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		cc.Driver,
		cc.User,
		cc.Pass,
		cc.Host,
		cc.Port,
		cc.Name,
	)
}

// Helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
