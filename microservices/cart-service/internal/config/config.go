package config

import "os"

type Config struct {
    ServiceName string
    ServicePort string
}

func Load() (*Config, error) {
    return &Config{ ServiceName: getEnv("SERVICE_NAME", "cart-service"), ServicePort: getEnv("SERVICE_PORT", "8083") }, nil
}

func getEnv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }