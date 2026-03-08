package config

import "os"

type Config struct {
    ServiceName string
    ServicePort string
}

func Load() (*Config, error) {
    return &Config{
        ServiceName: env("SERVICE_NAME", "product-service"),
        ServicePort: env("SERVICE_PORT", "8082"),
    }, nil
}

func env(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }