CREATE TABLE users (
    id    SERIAL PRIMARY KEY,
    email   TEXT UNIQUE    NOT NULL,
    password TEXT    NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);