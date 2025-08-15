package db

import (
  "fmt"
  "log"
  "os"
  "sync"
  "time"

  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
)

var (
  database *gorm.DB
  once     sync.Once
)

// SetDB initializes and stores the DB connection
func SetDB() {
  once.Do(func() {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    sslMode := os.Getenv("DB_SSLMODE") // optional, default "disable"

    if sslMode == "" {
      sslMode = "disable"
    }

    dsn := fmt.Sprintf(
      "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
      host, port, user, password, dbName, sslMode,
    )

    var err error
    database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
      Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
      log.Fatalf("failed to connect database: %v", err)
    }

    sqlDB, err := database.DB()
    if err != nil {
      log.Fatalf("failed to get sql.DB: %v", err)
    }

    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    log.Println("✅ Database connection established")
  })
}

// GetDB returns the existing DB instance
func GetDB() *gorm.DB {
  if database == nil {
    log.Fatal("❌ DB not initialized! Call SetDB() first.")
  }
  return database
}
