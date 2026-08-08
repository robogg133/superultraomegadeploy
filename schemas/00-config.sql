CREATE TABLE IF NOT EXISTS configs(
    key TEXT PRIMARY KEY,
    value JSONB NOT NULL,
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
    AFTER INSERT OR UPDATE ON configs
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
    AFTER DELETE ON configs
    FOR EACH ROW
    EXECUTE FUNCTION notify_config_delete();

CREATE OR REPLACE PROCEDURE set_config(
    p_key TEXT,
    p_value JSONB
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO configs (key, value)
    VALUES (p_key, p_value)
    ON CONFLICT (key)
    DO UPDATE SET
        value = p_value,
        updated_at = NOW();

    COMMIT;
END;
$$;
