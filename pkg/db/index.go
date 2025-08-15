package db

import (
  "fmt"
  "log"
  "sync"

  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
)

var (
  database *gorm.DB
  once     sync.Once
)

// SetDB initializes and stores the DB connection
func SetDB(dsn string) error {
  var err error
  database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
  })
  if err != nil {
    log.Fatalf("❌ Failed to open DB: %v", err)
    return err
  }
  fmt.Println("✅ Connected to PostgreSQL")

  // Auto migrate
  //db.AutoMigrate(&Service{}, &Task{}, &Execution{})
  return nil
}

// GetDB returns the existing DB instance
func GetDB() *gorm.DB {
  if database == nil {
    log.Fatal("❌ DB not initialized! Call SetDB() first.")
  }
  return database
}
