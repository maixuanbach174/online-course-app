-- Module queries

-- name: CreateModule :exec
INSERT INTO modules (id, course_id, title, order_index, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW());

-- name: UpdateModule :exec
UPDATE modules
SET title = $2,
    order_index = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteModule :exec
DELETE FROM modules WHERE id = $1;

-- name: GetModuleByID :one
SELECT id, course_id, title, order_index, created_at, updated_at
FROM modules
WHERE id = $1;

-- name: GetModulesByCourseID :many
SELECT id, course_id, title, order_index, created_at, updated_at
FROM modules
WHERE course_id = $1
ORDER BY order_index ASC;

-- name: ModuleExists :one
SELECT EXISTS(SELECT 1 FROM modules WHERE id = $1);

-- name: UpdateModuleOrder :exec
UPDATE modules
SET order_index = $2,
    updated_at = NOW()
WHERE id = $1;
