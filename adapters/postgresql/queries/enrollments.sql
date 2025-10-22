-- name: CreateEnrollment :exec
INSERT INTO enrollments (id, user_id, course_id, enrolled_at, started_at, completed_at, course_progress_percentage, course_progress_status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW());

-- name: CreateModuleProgress :exec
INSERT INTO module_progress (enrollment_id, module_id, progress_percentage, progress_status, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW());

-- name: CreateLessonProgress :exec
INSERT INTO lesson_progress (enrollment_id, lesson_id, progress_percentage, progress_status, exercise_score, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW());

-- name: UpdateEnrollment :exec
UPDATE enrollments
SET user_id = $2,
    course_id = $3,
    enrolled_at = $4,
    started_at = $5,
    completed_at = $6,
    course_progress_percentage = $7,
    course_progress_status = $8,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateModuleProgress :exec
UPDATE module_progress
SET progress_percentage = $3,
    progress_status = $4,
    updated_at = NOW()
WHERE enrollment_id = $1 AND module_id = $2;

-- name: UpdateLessonProgress :exec
UPDATE lesson_progress
SET progress_percentage = $3,
    progress_status = $4,
    exercise_score = $5,
    updated_at = NOW()
WHERE enrollment_id = $1 AND lesson_id = $2;

-- name: DeleteEnrollment :exec
DELETE FROM enrollments WHERE id = $1;

-- name: DeleteModuleProgressByEnrollmentID :exec
DELETE FROM module_progress WHERE enrollment_id = $1;

-- name: DeleteLessonProgressByEnrollmentID :exec
DELETE FROM lesson_progress WHERE enrollment_id = $1;

-- name: GetEnrollmentByID :one
SELECT id, user_id, course_id, enrolled_at, started_at, completed_at, course_progress_percentage, course_progress_status, created_at, updated_at
FROM enrollments
WHERE id = $1;

-- name: GetEnrollmentByUserAndCourse :one
SELECT id, user_id, course_id, enrolled_at, started_at, completed_at, course_progress_percentage, course_progress_status, created_at, updated_at
FROM enrollments
WHERE user_id = $1 AND course_id = $2;

-- name: GetModuleProgressByEnrollmentID :many
SELECT enrollment_id, module_id, progress_percentage, progress_status, created_at, updated_at
FROM module_progress
WHERE enrollment_id = $1
ORDER BY created_at ASC;

-- name: GetLessonProgressByEnrollmentID :many
SELECT enrollment_id, lesson_id, progress_percentage, progress_status, exercise_score, created_at, updated_at
FROM lesson_progress
WHERE enrollment_id = $1
ORDER BY created_at ASC;

-- name: GetAllEnrollments :many
SELECT id, user_id, course_id, enrolled_at, started_at, completed_at, course_progress_percentage, course_progress_status, created_at, updated_at
FROM enrollments
ORDER BY enrolled_at DESC;

-- name: GetEnrollmentsByUserID :many
SELECT id, user_id, course_id, enrolled_at, started_at, completed_at, course_progress_percentage, course_progress_status, created_at, updated_at
FROM enrollments
WHERE user_id = $1
ORDER BY enrolled_at DESC;
