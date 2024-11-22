-- name: CreateUrl :one
INSERT INTO
  urls (id, short_url, default_url)
VALUES
  (gen_random_uuid (), $1, $2)
RETURNING
  *;

-- name: GetUrl :one
SELECT
  *
FROM
  urls
WHERE
  short_url = $1;
