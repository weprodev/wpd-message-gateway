CREATE TABLE IF NOT EXISTS users
(
    id              UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    first_name      TEXT                              NOT NULL,
    last_name       TEXT                              NOT NULL,
    email           TEXT                              NOT NULL UNIQUE,
    password_hash   TEXT                              NOT NULL,
    email_verified  BOOLEAN                           NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ                       NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ                       NOT NULL DEFAULT NOW()
);
