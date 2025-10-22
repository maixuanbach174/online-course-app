-- name: CreateCourse :exec
INSERT INTO courses (id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW());

-- name: CreateCourseTag :exec
INSERT INTO course_tags (course_id, tag)
VALUES ($1, $2);

-- name: CreateModule :exec
INSERT INTO modules (id, course_id, title, order_index, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW());

-- name: CreateLesson :exec
INSERT INTO lessons (id, module_id, title, overview, content, video_id, order_index, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW());

-- name: CreateExercise :exec
INSERT INTO exercises (id, lesson_id, question, answers, correct_answer, order_index, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW());

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

-- name: DeleteCourseTags :exec
DELETE FROM course_tags WHERE course_id = $1;

-- name: DeleteModulesByCourseID :exec
DELETE FROM modules WHERE course_id = $1;

-- name: GetCourseByID :one
SELECT id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at
FROM courses
WHERE id = $1;

-- name: GetCourseTagsByCourseID :many
SELECT tag
FROM course_tags
WHERE course_id = $1
ORDER BY tag;

-- name: GetModulesByCourseID :many
SELECT id, course_id, title, order_index, created_at, updated_at
FROM modules
WHERE course_id = $1
ORDER BY order_index ASC;

-- name: GetLessonsByModuleID :many
SELECT id, module_id, title, overview, content, video_id, order_index, created_at, updated_at
FROM lessons
WHERE module_id = $1
ORDER BY order_index ASC;

-- name: GetExercisesByLessonID :many
SELECT id, lesson_id, question, answers, correct_answer, order_index, created_at, updated_at
FROM exercises
WHERE lesson_id = $1
ORDER BY order_index ASC;

-- name: GetAllCourses :many
SELECT id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at
FROM courses
ORDER BY created_at DESC;

-- name: GetCoursesByTeacherID :many
SELECT id, teacher_id, title, description, thumbnail, duration, domain, rating, level, created_at, updated_at
FROM courses
WHERE teacher_id = $1
ORDER BY created_at DESC;
