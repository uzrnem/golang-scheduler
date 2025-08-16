# Golang Scheduler Service (Echo + Postgres)

This microservice replaces heavy multiple Kubernetes CronJobs with a single lightweight Golang service that manages scheduled tasks.
Instead of running many CronJobs, this service schedules tasks in Postgres and executes them periodically based on their frequency and unit.

## Features

- REST APIs using Echo framework.
- PostgreSQL for persistent task storage.
- Scheduler that runs every 10 minutes to:
  - Fetch due tasks.
  - Insert execution logs.
  - Reschedule tasks.
  - Execute HTTP requests with headers and payloads.
- Supports GET, POST, PUT, DELETE, and OPTIONS HTTP methods.
- Secure per-service authentication using service_id and token.
- Dashboard to view and manage tasks.
- Stores request headers as JSONB in PostgreSQL.
- Pagination support for task executions.
- Dockerized setup for scheduler service, PostgreSQL, and pgAdmin.

## API Schema

1. Service Table

| Column      | Type      | Notes |
|-------------|-----------|-------|
| id          | uuid (PK) | Default generated |
| name        | string    |       |
| token       | string    | Auth key |
| created_at  | timestamp | Default now() |
| updated_at  | timestamp | Default now() |

2. Task Table

| Column        | Type              | Notes |
|---------------|-------------------|-------|
| id            | uuid (PK)         | Default generated |
| service_id    | uuid (FK)         | Ref to service.id |
| name          | string            |       |
| url           | string            | Target URL |
| header        | jsonb             | Request headers |
| payload       | text              | Request body |
| method        | string            | One of GET, POST, PUT, DELETE, OPTIONS |
| status        | string            | active, paused, disabled |
| scheduled_at  | timestamp         | Next run time |
| frequency     | integer           | > 0 |
| unit          | string            | hour or day |
| created_at    | timestamp         | Default now() |
| updated_at    | timestamp         | Default now() |

3. Execution Table

| Column      | Type      | Notes |
|-------------|-----------|-------|
| id          | uuid (PK) | Default generated |
| task_id     | uuid (FK) | Ref to task.id |
| status      | string    | started, success, failed |
| status_code | int       | HTTP status code |
| response    | text      | Response body |
| created_at  | timestamp | Default now() |

## Endpoints

### Create Task
POST /tasks

Headers:

  service_id: `uuid`

  token: `string`

Body:

```json
{
  "name": "Task Name",
  "url": "https://example.com",
  "header": { "Accept": "application/json" },
  "payload": "{}",
  "scheduled_at": "2025-08-14T12:00:00Z",  // or null
  "status": "active",  // active, paused, disabled
  "frequency": 1,
  "unit": "hour",
  "method": "GET"
}
```

### Update Task
POST /tasks/{taskId}

Headers:

  service_id: `uuid`

  token: `string`

Body:

```json
{
  "name": "Task Name",
  "url": "https://example.com",
  "header": { "Accept": "application/json" },
  "payload": "{}",
  "scheduled_at": "2025-08-14T12:00:00Z",  // or null
  "status": "active",  // active, paused, disabled
  "frequency": 1,
  "unit": "hour",
  "method": "GET"
}
```

### List Tasks
GET /tasks

Headers:

  service_id: `uuid`

  token: `string`

### Get Task by ID
GET /tasks/{taskId}

Headers:

  service_id: `uuid`

  token: `string`

### Get Task Executions (paginated)
GET /tasks/{taskId}/executions?pageNumber=0&count=10

Headers:

  service_id: `uuid`

  token: `string`

## Scheduler Behavior

- Runs every 10 minutes.
- Query: SELECT * FROM task WHERE scheduled_at <= now()
- For each due task:
  - Insert into execution table with status started.
  - Run the HTTP request with headers, payload, and method.
  - Update execution with status code, response, and status.
  - Reschedule task:
    scheduled_at = now() + (frequency * unit)

## Running the Project

### Development Mode (Live Code Mount)
Mounts your local source code into a Golang container for live changes.
```bash
docker compose --profile dev up
```

### Test/Production Mode (Build from Dockerfile)
Builds the Golang project into a lightweight container and runs it.
```bash
docker compose --profile test up --build
```

---

## Accessing the Services

- **Scheduler API** → http://localhost:8080 (Dashboard, check documents/schema.sql for Login details)
- **PostgreSQL** → localhost:5434 (user: `postgres`, password: `changeme`)  
- **pgAdmin** → http://localhost:8005 (login: `pgadmin4@pgadmin.org`, password: `admin`)

## Development Notes

- Migrations: Place SQL scripts in init.sql to auto-run at first container startup.
- DB persistence: Data is stored in ~/uzrnem/database/postgres_db.

## Environment Variables

| Variable            | Description                | Default |
|---------------------|----------------------------|---------|
| POSTGRES_USER       | Postgres username           | postgres |
| POSTGRES_PASSWORD   | Postgres password           | changeme |
| POSTGRES_HOST_AUTH_METHOD | Trust method         | trust |
| PGPORT              | Postgres container port     | 5434 |


---

## Folder Structure
```
.
├── cmd/
├── pkg/
│   ├── db/
│   ├── models/
│   ├── scheduler/
│   └── handlers/
├── dashboard/
│   └── index.html
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```# golang-scheduler
