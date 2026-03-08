package config

import (
    "fmt"
    "os"
    "time"
    "strings"
)

type AuthConfig struct {
    ServiceName      string
    ServicePort      string
    JWTSecret        string
    JWTAccessExpiry  time.Duration
    JWTRefreshExpiry time.Duration
    Database         DatabaseConfig
    Redis            RedisConfig
    Kafka            KafkaConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KafkaConfig struct {
    Brokers []string
    Topics  struct {
        UserRegistered string
        UserUpdated    string
        UserDeleted    string
        RoleAssigned   string
        RoleRevoked    string
    }
}

func LoadAuthConfig() (*AuthConfig, error) {
    return &AuthConfig{
        ServiceName:      getEnvOrDefault("SERVICE_NAME", "auth-service"),
        ServicePort:      getEnvOrDefault("SERVICE_PORT", "8081"),
        JWTSecret:        getEnvOrDefault("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
        JWTAccessExpiry:  getDurationEnvOrDefault("JWT_ACCESS_EXPIRY", 15*time.Minute),
        JWTRefreshExpiry: getDurationEnvOrDefault("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
        Database: DatabaseConfig{
            Host:     "",
            Port:     "",
            User:     "",
            Password: "",
            DBName:   getEnvOrDefault("DB_SQLITE_PATH", "/data/auth.db"),
            SSLMode:  "",
        },
        Redis: RedisConfig{
            Host:     "",
            Port:     "",
            Password: "",
            DB:       0,
        },
        Kafka: KafkaConfig{
            Brokers: getSliceEnvOrDefault("KAFKA_BROKERS", []string{"kafka:9092"}),
            Topics: struct {
                UserRegistered string
                UserUpdated    string
                UserDeleted    string
                RoleAssigned   string
                RoleRevoked    string
            }{
                UserRegistered: getEnvOrDefault("KAFKA_TOPIC_USER_REGISTERED", "user.registered"),
                UserUpdated:    getEnvOrDefault("KAFKA_TOPIC_USER_UPDATED", "user.updated"),
                UserDeleted:    getEnvOrDefault("KAFKA_TOPIC_USER_DELETED", "user.deleted"),
                RoleAssigned:   getEnvOrDefault("KAFKA_TOPIC_ROLE_ASSIGNED", "role.assigned"),
                RoleRevoked:    getEnvOrDefault("KAFKA_TOPIC_ROLE_REVOKED", "role.revoked"),
            },
        },
    }, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        var n int
        _, err := fmt.Sscanf(value, "%d", &n)
        if err == nil { return n }
    }
    return defaultValue
}

func getSliceEnvOrDefault(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        parts := strings.Split(value, ",")
        for i := range parts { parts[i] = strings.TrimSpace(parts[i]) }
        return parts
    }
    return defaultValue
}