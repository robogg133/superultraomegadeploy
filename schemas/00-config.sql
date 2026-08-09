CREATE SCHEMA IF NOT EXISTS configs;

CREATE TABLE IF NOT EXISTS configs.configs(
    key TEXT PRIMARY KEY,
    value JSONB NOT NULL,
    version INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION notify_config_change()
RETURNS TRIGGER AS $$
DECLARE
    payload jsonb;
BEGIN
    payload = jsonb_build_object(
        'key', NEW.key,
        'value', NEW.value,
        'operation', TG_OP,
        'updated_at', NEW.updated_at
    );

    PERFORM pg_notify('config_changes', payload::text);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER config_notify_trigger
    AFTER INSERT OR UPDATE ON configs.configs
    FOR EACH ROW
    EXECUTE FUNCTION notify_config_change();

CREATE OR REPLACE FUNCTION notify_config_delete()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('config_changes',
        jsonb_build_object(
            'key', OLD.key,
            'operation', 'DELETE',
            'updated_at', NOW()
        )::text
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER config_delete_trigger
    AFTER DELETE ON configs.configs
    FOR EACH ROW
    EXECUTE FUNCTION notify_config_delete();

CREATE OR REPLACE PROCEDURE set_config(
    p_key TEXT,
    p_value JSONB
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO configs.configs (key, value)
    VALUES (p_key, p_value)
    ON CONFLICT (key)
    DO UPDATE SET
        value = p_value,
        version = configs.configs.version + 1,
        updated_at = NOW();

    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE set_server_config(
    p_user_id UUID,
    p_key TEXT,
    p_value JSONB
)
LANGUAGE plpgsql AS $$
DECLARE
    can_change BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1 FROM users.permissions
        WHERE user_id = p_user_id
          AND permission = 'canChangeServerSettings'
    ) INTO can_change;

    IF can_change IS NOT TRUE THEN
        RAISE EXCEPTION 'forbidden'
            USING ERRCODE = '42501';
    END IF;

    CALL set_config(p_key, p_value);
END;
$$;
