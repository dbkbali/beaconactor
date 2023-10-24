package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DATABASE_URL", "test_database_url")
	os.Setenv("BEACON_NODE_URL", "test_beacon_node_url")
	os.Setenv("DB_USER", "test_db_user")
	os.Setenv("DB_PASSWORD", "test_db_password")
	os.Setenv("DB_HOST", "test_db_host")
	os.Setenv("DB_PORT", "test_db_port")
	os.Setenv("DB_NAME", "test_db_name")

	// Call the New function
	config := New()

	// Check that the values are set correctly
	if config.DatabaseUrl != "test_database_url" {
		t.Errorf("expected DatabaseUrl to be 'test_database_url', got '%s'", config.DatabaseUrl)
	}
	if config.BeaconNodeUrl != "test_beacon_node_url" {
		t.Errorf("expected BeaconNodeUrl to be 'test_beacon_node_url', got '%s'", config.BeaconNodeUrl)
	}
	if config.DbUser != "test_db_user" {
		t.Errorf("expected DbUser to be 'test_db_user', got '%s'", config.DbUser)
	}
	if config.DbPassword != "test_db_password" {
		t.Errorf("expected DbPassword to be 'test_db_password', got '%s'", config.DbPassword)
	}
	if config.DbHost != "test_db_host" {
		t.Errorf("expected DbHost to be 'test_db_host', got '%s'", config.DbHost)
	}
	if config.DbPort != "test_db_port" {
		t.Errorf("expected DbPort to be 'test_db_port', got '%s'", config.DbPort)
	}
	if config.DbName != "test_db_name" {
		t.Errorf("expected DbName to be 'test_db_name', got '%s'", config.DbName)
	}
}

func TestGetEnv(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TEST_KEY", "test_value")

	// Call the getEnv function with an existing key
	value := getEnv("TEST_KEY", "fallback_value")
	if value != "test_value" {
		t.Errorf("expected value to be 'test_value', got '%s'", value)
	}

	// Call the getEnv function with a non-existing key
	value = getEnv("NON_EXISTING_KEY", "fallback_value")
	if value != "fallback_value" {
		t.Errorf("expected value to be 'fallback_value', got '%s'", value)
	}
}
