-- name: GetVisitCount :one
SELECT COUNT(*) FROM visits;

-- name: ListVisits :many
SELECT * FROM visits ORDER BY id LIMIT $1 OFFSET $2;

-- name: CreateVisit :one
INSERT INTO visits (link_id, ip, user_agent, referer, "status") VALUES ($1, $2, $3, $4, $5) RETURNING *;
