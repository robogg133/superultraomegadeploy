CREATE SCHEMA IF NOT EXISTS project;

CREATE TABLE IF NOT EXISTS project.project(
    org_id UUID REFERENCES org.org(org_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    project_id UUID PRIMARY KEY DEFAULT uuidv7() , 

    project_slug TEXT NOT NULL,
    project_name TEXT NOT NULL,


    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, project_slug)
);

CREATE TABLE IF NOT EXISTS project.permission (
    org_member_id UUID PRIMARY KEY REFERENCES org.members(org_member_id) ON DELETE CASCADE ON UPDATE CASCADE,

    can_read_project BOOL NOT NULL DEFAULT TRUE,
    can_write_project BOOL NOT NULL DEFAULT TRUE,
    can_rdwr_project_secrets BOOL NOT NULL DEFAULT TRUE
    
);

CREATE TABLE IF NOT EXISTS project.secret (
    project_id UUID PRIMARY KEY REFERENCES project.project(project_id) ON DELETE CASCADE ON UPDATE CASCADE,
    secrets_env_file TEXT NOT NULL DEFAULT ''
);
