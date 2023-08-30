CREATE TABLE IF NOT EXISTS segments (
    id SERIAL PRIMARY KEY NOT NULL,
    slug VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_segments (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users,
    segment_id INTEGER NOT NULL REFERENCES segments,
    created_at TIMESTAMP NOT NULL,
    ttl TIMESTAMP,
    UNIQUE (user_id, segment_id)
);
CREATE TABLE IF NOT EXISTS operations (
                                          id BIGSERIAL NOT NULL PRIMARY KEY,
                                          user_id BIGINT NOT NULL,
                                          segment_id BIGINT NOT NULL,
                                          operation VARCHAR(16) NOT NULL,
                                          time TIMESTAMP NOT NULL
);

CREATE OR REPLACE FUNCTION process_operations_log() RETURNS TRIGGER AS $operations$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO operations (user_id, segment_id, operation, time) VALUES (OLD.user_id, OLD.segment_id, 'delete', now());
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO operations (user_id, segment_id, operation, time) VALUES (NEW.user_id, NEW.segment_id, 'add', now());
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$operations$ LANGUAGE plpgsql;

CREATE TRIGGER operations_log
    AFTER INSERT OR DELETE ON users_segments
    FOR EACH ROW EXECUTE PROCEDURE process_operations_log();

INSERT INTO users (name)
VALUES
    ('admin'),
    ('root'),
    ('user');

INSERT INTO segments (slug)
VALUES
    ('AVITO_VOICE_MESSAGES'),
    ('AVITO_SALE_30'),
    ('AVITO_SALE_70');

INSERT INTO users_segments (user_id, segment_id, created_at)
VALUES
    (1, 1, now()),
    (1, 2, now()),
    (2, 3, now()),
    (3, 2, now()),
    (3, 3, now());