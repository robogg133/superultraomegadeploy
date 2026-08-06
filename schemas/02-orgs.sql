CREATE SCHEMA IF NOT EXISTS org;

CREATE TABLE IF NOT EXISTS org.org (
    org_id UUID PRIMARY KEY DEFAULT uuidv7(),

    org_slug TEXT NOT NULL,
    org_name TEXT NOT NULL,

    created_by UUID REFERENCES users.user(user_id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(created_by, org_name)
);

CREATE TABLE IF NOT EXISTS org.members (
    org_id UUID REFERENCES org.org(org_id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id UUID REFERENCES users.user(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    org_member_id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Permssions
    can_create_projects BOOL NOT NULL DEFAULT TRUE,
    can_modify_org BOOL NOT NULL DEFAULT TRUE,



    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_user_id ON org.members(user_id);
CREATE INDEX IF NOT EXISTS idx_org_id_user_id ON org.members(org_id, user_id);
