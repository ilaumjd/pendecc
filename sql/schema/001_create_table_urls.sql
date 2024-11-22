-- +goose Up
CREATE TABLE urls (
  id UUID PRIMARY KEY,
  short_url TEXT NOT NULL UNIQUE,
  default_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE urls;
