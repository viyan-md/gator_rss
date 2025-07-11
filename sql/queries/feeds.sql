-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :one 
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeeds :many
SELECT
  f.name   AS feed_name,
  f.url    AS feed_url,
  u.name   AS user_name
FROM feeds AS f
LEFT JOIN users AS u
  ON f.user_id = u.id;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one 
SELECT *
FROM feeds
ORDER BY
  last_fetched_at NULLS FIRST,  
  created_at ASC      
LIMIT 1;