package main

import (
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "go-echo-scheduler/pkg/scheduler"
  "go-echo-scheduler/pkg/handlers"
  "go-echo-scheduler/pkg/db"
)

func main() {
  db.SetDB() // Initialize database connection

  // Start scheduler (runs every 10 mins)
  go scheduler.StartScheduler()

  // Echo instance
  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // Serve static dashboard
  e.Static("/dashboard", "dashboard")
  e.GET("/", func(c echo.Context) error {
    return c.File("dashboard/index.html")
  })

  // API routes
  e.POST("/tasks", handlers.CreateTask)
  e.GET("/tasks", handlers.GetTasks)
  e.GET("/tasks/:taskId", handlers.GetTaskByID)
  e.PUT("/tasks/:taskId", handlers.UpdateTaskByID)
  e.GET("/tasks/:taskId/executions", handlers.GetTaskExecutions)

  // Start server
  e.Logger.Fatal(e.Start(":8080"))
}