-- Lesson queries

-- name: CreateLesson :exec
INSERT INTO lessons (id, module_id, title, overview, content, video_id, duration, order_index, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW());

-- name: UpdateLesson :exec
UPDATE lessons
SET title = $2,
    overview = $3,
    content = $4,
    video_id = $5,
    duration = $6,
    order_index = $7,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteLesson :exec
DELETE FROM lessons WHERE id = $1;

-- name: GetLessonByID :one
SELECT id, module_id, title, overview, content, video_id, duration, order_index, created_at, updated_at
FROM lessons
WHERE id = $1;

-- name: GetLessonsByModuleID :many
SELECT id, module_id, title, overview, content, video_id, duration, order_index, created_at, updated_at
FROM lessons
WHERE module_id = $1
ORDER BY order_index ASC;

-- name: LessonExists :one
SELECT EXISTS(SELECT 1 FROM lessons WHERE id = $1);

-- name: UpdateLessonOrder :exec
UPDATE lessons
SET order_index = $2,
    updated_at = NOW()
WHERE id = $1;
