-- name: CreateFeedFollow :one
WITH inserted AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT
  inserted.*,
  u.name AS user_name,
  f.name AS feed_name
FROM inserted
JOIN users u ON inserted.user_id = u.id
JOIN feeds f ON inserted.feed_id = f.id;

-- name: GetFeedFollowsForUser :many 
SELECT f.name AS feed_name 
FROM feed_follows ff
JOIN feeds f ON ff.feed_id = f.id
JOIN users u ON ff.user_id = u.id
WHERE u.id = $1;

-- name: UnfollowFeed :one
DELETE FROM feed_follows ff
USING users u, feeds f
WHERE
    ff.user_id = u.id
    AND ff.feed_id = f.id
    AND u.id      = $1
    AND f.url    = $2
RETURNING ff.*;
