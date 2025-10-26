-- Exercise queries

-- name: CreateExercise :exec
INSERT INTO exercises (id, lesson_id, question, answers, correct_answer, order_index, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW());

-- name: UpdateExercise :exec
UPDATE exercises
SET question = $2,
    answers = $3,
    correct_answer = $4,
    order_index = $5,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteExercise :exec
DELETE FROM exercises WHERE id = $1;

-- name: GetExerciseByID :one
SELECT id, lesson_id, question, answers, correct_answer, order_index, created_at, updated_at
FROM exercises
WHERE id = $1;

-- name: GetExercisesByLessonID :many
SELECT id, lesson_id, question, answers, correct_answer, order_index, created_at, updated_at
FROM exercises
WHERE lesson_id = $1
ORDER BY order_index ASC;

-- name: ExerciseExists :one
SELECT EXISTS(SELECT 1 FROM exercises WHERE id = $1);

-- name: UpdateExerciseOrder :exec
UPDATE exercises
SET order_index = $2,
    updated_at = NOW()
WHERE id = $1;
