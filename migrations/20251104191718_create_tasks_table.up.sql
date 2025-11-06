CREATE TABLE tasks (
    id     SERIAL PRIMARY KEY,
    text   TEXT    NOT NULL,
    status TEXT    NOT NULL DEFAULT 'new'
);