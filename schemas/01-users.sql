CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.user (
    user_id UUID PRIMARY KEY DEFAULT uuidv7(),

    user_email TEXT NOT NULL UNIQUE,
    user_first_name TEXT NOT NULL,
    user_last_name TEXT NOT NULL,

    user_password TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_user_email ON users.user(user_email);

-- This table needs a automatic cleanup
CREATE TABLE users.sessions (
    session_id TEXT PRIMARY KEY,
    user_id UUID REFERENCES users.user(user_id) ON DELETE CASCADE ON UPDATE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON users.sessions(user_id);

CREATE TABLE IF NOT EXISTS users.permissions (
    user_id UUID REFERENCES users.user(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    permission TEXT NOT NULL,

    PRIMARY KEY (user_id, permission)
);
