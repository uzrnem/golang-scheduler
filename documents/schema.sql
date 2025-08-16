CREATE EXTENSION IF NOT EXISTS hstore;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Service table
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Task table
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL CHECK (method IN ('GET', 'POST', 'DELETE', 'PUT', 'OPTIONS')),
    header jsonb,
    payload TEXT,
    scheduled_at TIMESTAMP NOT NULL,
    frequency INTEGER NOT NULL CHECK (frequency > 0),
    unit TEXT NOT NULL CHECK (unit IN ('hour', 'day')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
ALTER TABLE tasks
    ADD COLUMN status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'disabled'));
-- ALTER TABLE tasks ALTER COLUMN header TYPE jsonb USING header::jsonb;


-- Execution table
CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    status_code INTEGER,
    response TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_tasks_service_id ON tasks(service_id);
CREATE INDEX idx_tasks_scheduled_at ON tasks(scheduled_at);
CREATE INDEX idx_executions_task_id ON executions(task_id);
CREATE INDEX idx_executions_created_at ON executions(created_at DESC);

-- ALTER TABLE executions ALTER COLUMN id SET DEFAULT uuid_generate_v4();


-- # Seed
insert into services (id, name, token) 
values ('41a0f37d-706d-4f61-b115-e98b817ed360', 'int svc', 'secret-token');