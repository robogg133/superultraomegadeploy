CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.user (
    user_id UUID PRIMARY KEY DEFAULT uuidv7(),

    user_email TEXT NOT NULL UNIQUE,
    user_first_name TEXT NOT NULL,
    user_last_name TEXT NOT NULL,
    
    user_password TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);
CREATE INDEX IF NOT EXISTS idx_user_email ON users.user(user_email);

CREATE TABLE IF NOT EXISTS users.permission (
    user_id UUID PRIMARY KEY REFERENCES users.user(user_id) ON DELETE CASCADE ON UPDATE CASCADE,

    can_create_org  BOOL NOT NULL DEFAULT FALSE,
    can_invite_user BOOL NOT NULL DEFAULT FALSE,
    can_change_server_settings  BOOL NOT NULL DEFAULT FALSE
);

