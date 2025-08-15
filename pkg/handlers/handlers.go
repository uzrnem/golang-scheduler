package handlers

import (
  "fmt"
  "net/http"
  "time"
  "strings"
  "strconv"

  "github.com/google/uuid"
  "github.com/labstack/echo/v4"
  "go-echo-scheduler/pkg/db"
  "go-echo-scheduler/pkg/models"
)

var validMethods = map[string]bool{
  "GET":     true,
  "POST":    true,
  "DELETE":  true,
  "PUT":     true,
  "OPTIONS": true,
}

// Middleware to validate service_id & token
func validateService(c echo.Context) (*models.Service, error) {
  serviceIDStr := c.Request().Header.Get("service_id")
  token := c.Request().Header.Get("token")

  serviceID, err := uuid.Parse(serviceIDStr)
  if err != nil {
    return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid service_id")
  }
  fmt.Println("Validating service with ID:", serviceID)

  var service models.Service
  if err := db.GetDB().First(&service, "id = ? AND token = ?", serviceID, token).Error; err != nil {
    return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid service credentials")
  }
  fmt.Println("Validated service:", service.ID)
  return &service, nil
}

// POST /tasks
func CreateTask(c echo.Context) error {
  service, err := validateService(c)
  if err != nil {
    return err
  }

  var req models.Task
  if err := c.Bind(&req); err != nil {
    fmt.Println("Error creating task:", err)
    return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
  }
  fmt.Println("Creating task for req:", req)

  if req.Frequency <= 0 {
    return echo.NewHTTPError(http.StatusBadRequest, "frequency must be greater than zero")
  }
  if req.Unit != "hour" && req.Unit != "day" {
    return echo.NewHTTPError(http.StatusBadRequest, "unit must be 'hour' or 'day'")
  }
  req.Method = strings.ToUpper(req.Method)
  if !validMethods[req.Method] {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid method, allowed: GET, POST, DELETE, PUT, OPTIONS")
  }

  req.ID = uuid.New()
  req.ServiceID = service.ID
  req.CreatedAt = time.Now()
  req.UpdatedAt = time.Now()

  // If ScheduledAt is null, set based on frequency
  if req.ScheduledAt.IsZero() {
    if req.Unit == "hour" {
      req.ScheduledAt = time.Now().Add(time.Duration(req.Frequency) * time.Hour)
    } else {
      req.ScheduledAt = time.Now().Add(time.Duration(req.Frequency) * 24 * time.Hour)
    }
  }

  if err := db.GetDB().Create(&req).Error; err != nil {
    return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
  }
  return c.JSON(http.StatusCreated, req)
}

// GET /tasks
func GetTasks(c echo.Context) error {
  service, err := validateService(c)
  if err != nil {
    return err
  }

  var tasks []models.Task
  if err := db.GetDB().Where("service_id = ?", service.ID).Find(&tasks).Error; err != nil {
    return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
  }
  return c.JSON(http.StatusOK, tasks)
}

// GET /tasks/:taskId
func GetTaskByID(c echo.Context) error {
  service, err := validateService(c)
  if err != nil {
    return err
  }

  taskID, err := uuid.Parse(c.Param("taskId"))
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid taskId")
  }

  var task models.Task
  if err := db.GetDB().First(&task, "id = ? AND service_id = ?", taskID, service.ID).Error; err != nil {
    return echo.NewHTTPError(http.StatusNotFound, "task not found")
  }
  return c.JSON(http.StatusOK, task)
}

// GET /tasks/:taskId/executions
func GetTaskExecutions(c echo.Context) error {
  service, err := validateService(c)
  if err != nil {
    return err
  }

  taskID, err := uuid.Parse(c.Param("taskId"))
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid taskId")
  }

  // Ensure task belongs to service
  var count int64
  db.GetDB().Model(&models.Task{}).Where("id = ? AND service_id = ?", taskID, service.ID).Count(&count)
  if count == 0 {
    return echo.NewHTTPError(http.StatusNotFound, "task not found")
  }

  // Pagination params
  pageNumber, err := strconv.Atoi(c.QueryParam("pageNumber"))
  if err != nil || pageNumber < 0 {
    pageNumber = 0
  }
  limit, err := strconv.Atoi(c.QueryParam("count"))
  if err != nil || limit <= 0 {
    limit = 10
  }
  offset := pageNumber * limit

  // Total executions count
  var totalCount int64
  db.GetDB().Model(&models.Execution{}).Where("task_id = ?", taskID).Count(&totalCount)

  // Fetch executions with pagination
  var executions []models.Execution
  if err := db.GetDB().Where("task_id = ?", taskID).
    Order("created_at DESC").
    Limit(limit).
    Offset(offset).
    Find(&executions).Error; err != nil {
    return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
  }

  return c.JSON(http.StatusOK, map[string]interface{}{
    "total_count": totalCount,
    "page_number": pageNumber,
    "count":       limit,
    "records":     executions,
  })
}
