package scheduler

import (
  "io"
  "net/http"
  "context"
  "log"
  "sync"
  "time"
  "strings"

  "github.com/google/uuid"
  "go-echo-scheduler/pkg/db"
  "go-echo-scheduler/pkg/models"
)

func StartScheduler() {
  ticker := time.NewTicker(10 * time.Minute)
  defer ticker.Stop()

  for {
    select {
    case <-ticker.C:
      runScheduler()
    }
  }
}

func runScheduler() {
  log.Println("Scheduler triggered:", time.Now())

  var tasks []models.Task
  if err := db.GetDB().Where("scheduled_at <= ?", time.Now()).Find(&tasks).Error; err != nil {
    log.Println("Error fetching tasks:", err)
    return
  }

  if len(tasks) == 0 {
    log.Println("No due tasks.")
    return
  }

  log.Printf("Found %d due tasks\n", len(tasks))

  // Fixed pool size
  poolSize := 4
  taskChan := make(chan models.Task)
  var wg sync.WaitGroup

  // Start workers
  for i := 0; i < poolSize; i++ {
    wg.Add(1)
    go func(workerID int) {
      defer wg.Done()
      for task := range taskChan {
        log.Printf("Worker %d started task %s\n", workerID, task.ID)
        processTask(task)
        log.Printf("Worker %d finished task %s\n", workerID, task.ID)
        }
    }(i + 1)
  }

  // Feed tasks into the pool
  go func() {
    for _, t := range tasks {
      taskChan <- t
    }
    close(taskChan)
  }()

  wg.Wait()
}

// processTask handles the DB updates, execution record creation, and HTTP call
func processTask(task models.Task) {
  // 1. Insert into execution table with "started" status
  exec := models.Execution{
    ID:        uuid.New(),
    TaskID:    task.ID,
    Status:    "started",
    CreatedAt: time.Now(),
  }
  if err := db.GetDB().Create(&exec).Error; err != nil {
    log.Println("Error creating execution:", err)
    return
  }

  // 2. Make HTTP request to task.URL with headers & payload
  statusCode, respBody := executeHTTPRequest(task)

  status := "completed"
  if statusCode >= 400 || statusCode == 0 {
    log.Printf("Task %s failed with status code %d\n", task.ID, status)
    status = "failed"
  }

  // 3. Update execution with result
  if err := db.GetDB().Model(&exec).Updates(map[string]interface{}{
    "status":      status,
    "status_code": statusCode,
    "response":    respBody,
  }).Error; err != nil {
    log.Println("Error updating execution:", err)
  }

  // 4. Update scheduled_at for next run
  nextTime := task.ScheduledAt
  if task.Unit == "hour" {
    nextTime = nextTime.Add(time.Duration(task.Frequency) * time.Hour)
  } else if task.Unit == "day" {
    nextTime = nextTime.Add(time.Duration(task.Frequency*24) * time.Hour)
  }

  if err := db.GetDB().Model(&task).Update("scheduled_at", nextTime).Error; err != nil {
    log.Println("Error updating next scheduled time:", err)
  }
}

func executeHTTPRequest(task models.Task) (int, string) {
  ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  defer cancel()

  req, err := http.NewRequestWithContext(ctx, task.Method, task.URL, strings.NewReader(task.Payload))
  if err != nil {
    log.Println("HTTP request creation error:", err)
    return 0, err.Error()
  }

  // Add headers from JSONB
  for k, v := range task.Header {
    req.Header.Set(k, v.(string))
  }

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    statusCode := 0
    if resp != nil {
      statusCode = resp.StatusCode
    }
    log.Println("HTTP request error:", err)
    return statusCode, err.Error()
  }
  defer resp.Body.Close()

  bodyBytes, _ := io.ReadAll(resp.Body)
  return resp.StatusCode, string(bodyBytes)
}