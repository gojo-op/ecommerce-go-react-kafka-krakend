package config

import "os"

type Config struct {
    ServiceName string
    ServicePort string
}

func Load() (*Config, error) {
    return &Config{ ServiceName: env("SERVICE_NAME","order-service"), ServicePort: env("SERVICE_PORT","8084") }, nil
}

func env(k,d string) string { if v := os.Getenv(k); v != "" { return v }; return d }