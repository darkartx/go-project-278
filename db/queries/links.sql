-- name: ListLinks :many
SELECT * FROM links ORDER BY id;

-- name: CreateLink :one
INSERT INTO links (original_url, short_name) VALUES ($1, $2) RETURNING *;

-- name: GetLink :one
SELECT * FROM links WHERE id = $1;

-- name: UpdateLink :one
UPDATE links SET original_url = $1, short_name = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING *;

-- name: DeleteLink :exec
DELETE FROM links WHERE id = $1;
