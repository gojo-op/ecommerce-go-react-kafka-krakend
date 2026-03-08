package db

import (
    "os"
    "path/filepath"
    sqlite "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
    path := os.Getenv("DB_SQLITE_PATH")
    if path == "" { path = "/data/auth.db" }
    _ = os.MkdirAll(filepath.Dir(path), 0o755)
    db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
    if err == nil { return db, nil }
    // Fallback to local file if volume not writable
    localPath := "/tmp/auth.db"
    _ = os.MkdirAll(filepath.Dir(localPath), 0o755)
    return gorm.Open(sqlite.Open(localPath), &gorm.Config{})
}