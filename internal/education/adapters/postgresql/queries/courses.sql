-- Course queries

-- name: CreateCourse :exec
INSERT INTO courses (id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW());

-- name: UpdateCourse :exec
UPDATE courses
SET teacher_id = $2,
    title = $3,
    description = $4,
    thumbnail = $5,
    duration = $6,
    domain = $7,
    rating = $8,
    level = $9,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1;

-- name: GetCourseByID :one
SELECT id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at
FROM courses
WHERE id = $1;

-- name: GetAllCourses :many
SELECT id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at
FROM courses
ORDER BY created_at DESC;

-- name: GetCoursesByTeacherID :many
SELECT id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at
FROM courses
WHERE teacher_id = $1
ORDER BY created_at DESC;

-- name: CourseExists :one
SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1);

-- Course Tags queries

-- name: CreateCourseTag :exec
INSERT INTO course_tags (course_id, tag)
VALUES ($1, $2)
ON CONFLICT (course_id, tag) DO NOTHING;

-- name: DeleteCourseTag :exec
DELETE FROM course_tags WHERE course_id = $1 AND tag = $2;

-- name: DeleteAllCourseTags :exec
DELETE FROM course_tags WHERE course_id = $1;

-- name: GetCourseTagsByCourseID :many
SELECT tag
FROM course_tags
WHERE course_id = $1
ORDER BY tag;
