package config

import (
    "os"
    "fmt"
    "strings"
)

type Config struct {
    Environment string
    Database struct {
        Host string
        Port string
        User string
        Password string
        Name string
        SSLMode string
    }
    Redis struct {
        Host string
        Port string
        Password string
        DB int
    }
    Kafka struct{
        Brokers []string
    }
}

func LoadConfig() (*Config, error) {
    c := &Config{}
    c.Environment = getenv("ENVIRONMENT", "development")
    c.Database.Host = getenv("DB_HOST", "postgres")
    c.Database.Port = getenv("DB_PORT", "5432")
    c.Database.User = getenv("DB_USER", "postgres")
    c.Database.Password = getenv("DB_PASSWORD", "password")
    c.Database.Name = getenv("DB_NAME", "demo_app")
    c.Database.SSLMode = getenv("DB_SSL_MODE", "disable")
    c.Redis.Host = getenv("REDIS_HOST", "redis")
    c.Redis.Port = getenv("REDIS_PORT", "6379")
    c.Redis.Password = getenv("REDIS_PASSWORD", "")
    c.Redis.DB = intFromEnv("REDIS_DB", 1)
    c.Kafka.Brokers = []string{ getenv("KAFKA_BROKERS", "kafka:9092") }
    return c, nil
}

func getenv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }
func intFromEnv(k string, d int) int { if v := os.Getenv(k); v != "" { if n, err := ParseInt(v); err == nil { return n } }; return d }

func ParseInt(s string) (int, error) {
    var n int
    _, err := fmt.Sscanf(s, "%d", &n)
    return n, err
}

func ParseSlice(s string) []string {
    if s == "" { return []string{} }
    parts := strings.Split(s, ",")
    for i := range parts { parts[i] = strings.TrimSpace(parts[i]) }
    return parts
}