package models

import (
  "time"

  "github.com/google/uuid"
  "gorm.io/datatypes"
)

// Service table
type Service struct {
  ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
  Name      string    `json:"name"`
  Token     string    `json:"token"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

// Task table
type Task struct {
  ID          uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
  ServiceID   uuid.UUID         `gorm:"type:uuid;index" json:"service_id"`
  Name        string            `json:"name"`
  URL         string            `json:"url"`
  Method      string            `json:"method"` // Added method column
  Header      datatypes.JSONMap `gorm:"type:hstore" json:"header"`
  Payload     string            `json:"payload"`
  ScheduledAt time.Time         `json:"scheduled_at"`
  Frequency   int               `json:"frequency"`
  Unit        string            `json:"unit"` // "hour" or "day"
  Status      string            `json:"status"` // "active", "paused", "disabled"
  CreatedAt   time.Time         `json:"created_at"`
  UpdatedAt   time.Time         `json:"updated_at"`
}

// Execution table
type Execution struct {
  ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
  TaskID     uuid.UUID `gorm:"type:uuid;index" json:"task_id"`
  Status     string    `json:"status"`
  StatusCode int       `json:"status_code"`
  Response   string    `json:"response"`
  CreatedAt  time.Time `json:"created_at"`
}
