-- name: CreateFeedFollows :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT inserted.*, feeds.name as feed_name, users.name as user_name
FROM inserted
INNER JOIN feeds ON inserted.feed_id = feeds.id
INNER JOIN users ON inserted.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.id as follow_id, users.name as user_name, feeds.name as feed_name FROM feed_follows
LEFT JOIN users on feed_follows.user_id = users.id
LEFT JOIN feeds on feed_follows.feed_id = feeds.id
WHERE users.name = $1;