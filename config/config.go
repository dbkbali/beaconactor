package config

import "os"

type Config struct {
	DatabaseUrl   string
	BeaconNodeUrl string
	DbUser        string
	DbPassword    string
	DbHost        string
	DbPort        string
	DbName        string
}

func New() *Config {
	return &Config{
		DatabaseUrl:   getEnv("DATABASE_URL", ""),
		BeaconNodeUrl: getEnv("BEACON_NODE_URL", ""),
		DbUser:        getEnv("DB_USER", ""),
		DbPassword:    getEnv("DB_PASSWORD", ""),
		DbHost:        getEnv("DB_HOST", "localhost"),
		DbPort:        getEnv("DB_PORT", "5432"),
		DbName:        getEnv("DB_NAME", "beaconvalidators"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
