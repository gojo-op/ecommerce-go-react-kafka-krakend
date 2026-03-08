package database

import (
    "os"
    "path/filepath"
    sqlite "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    cfgpkg "github.com/your-org/microservices/shared/config"
)

type DB struct { DB *gorm.DB }

func NewDatabase(cfg *cfgpkg.Config) (*DB, error) {
    path := os.Getenv("DB_SQLITE_PATH")
    if path == "" {
        name := cfg.Database.Name
        if name == "" { name = "database" }
        path = filepath.Join("/data", name+".db")
    }
    _ = os.MkdirAll(filepath.Dir(path), 0o755)
    g, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
    if err != nil { return nil, err }
    return &DB{ DB: g }, nil
}

func (d *DB) AutoMigrate(models ...interface{}) error { return d.DB.AutoMigrate(models...) }
func (d *DB) Close() error { return nil }
