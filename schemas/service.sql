CREATE SCHEMA IF NOT EXISTS service;

CREATE TABLE IF NOT EXISTS service.service (
    project_id BIGINT REFERENCES projects.project(project_id), 

    service_id UUID PRIMARY KEY DEFAULT uuidv7(),
    service_slug TEXT NOT NULL,
    service_name TEXT NOT NULL,


    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(project_id, service_slug)
);

