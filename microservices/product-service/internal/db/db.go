package db

import (
    "os"
    "path/filepath"
    sqlite "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
    path := os.Getenv("DB_SQLITE_PATH")
    if path == "" { path = "/data/product.db" }
    _ = os.MkdirAll(filepath.Dir(path), 0o755)
    return gorm.Open(sqlite.Open(path), &gorm.Config{})
}