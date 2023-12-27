CREATE TABLE IF NOT EXISTS urls (
    id bigserial PRIMARY KEY,
    shortKey text NOT NULL,
    longURL text NOT NULL
);