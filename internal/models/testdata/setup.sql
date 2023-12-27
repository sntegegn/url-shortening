CREATE TABLE IF NOT EXISTS urls (
    id bigserial PRIMARY KEY,
    shortKey text NOT NULL,
    longURL text NOT NULL
);

INSERT INTO urls (shortKey, longURL) VALUES (
    'mUV4W2',
    'http://example.com'
)