-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'student', 'teacher')),
    profile TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Courses table
CREATE TABLE IF NOT EXISTS courses (
    id VARCHAR(255) PRIMARY KEY,
    teacher_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    thumbnail TEXT,
    duration INT NOT NULL DEFAULT 0,
    domain VARCHAR(100) NOT NULL,
    rating DECIMAL(3, 2) DEFAULT 0.0,
    level VARCHAR(50) NOT NULL CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_courses_teacher_id ON courses(teacher_id);
CREATE INDEX idx_courses_domain ON courses(domain);
CREATE INDEX idx_courses_level ON courses(level);

-- Course tags (many-to-many through junction table)
CREATE TABLE IF NOT EXISTS course_tags (
    course_id VARCHAR(255) NOT NULL,
    tag VARCHAR(100) NOT NULL,
    PRIMARY KEY (course_id, tag),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE INDEX idx_course_tags_tag ON course_tags(tag);

-- Modules table
CREATE TABLE IF NOT EXISTS modules (
    id VARCHAR(255) PRIMARY KEY,
    course_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    order_index INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE INDEX idx_modules_course_id ON modules(course_id);
CREATE INDEX idx_modules_order ON modules(course_id, order_index);

-- Lessons table
CREATE TABLE IF NOT EXISTS lessons (
    id VARCHAR(255) PRIMARY KEY,
    module_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    overview TEXT,
    content TEXT,
    video_id VARCHAR(255),
    order_index INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE
);

CREATE INDEX idx_lessons_module_id ON lessons(module_id);
CREATE INDEX idx_lessons_order ON lessons(module_id, order_index);

-- Exercises table
CREATE TABLE IF NOT EXISTS exercises (
    id VARCHAR(255) PRIMARY KEY,
    lesson_id VARCHAR(255) NOT NULL,
    question TEXT NOT NULL,
    answers TEXT[] NOT NULL,
    correct_answer TEXT NOT NULL,
    order_index INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);

CREATE INDEX idx_exercises_lesson_id ON exercises(lesson_id);

-- Enrollments table
CREATE TABLE IF NOT EXISTS enrollments (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    course_id VARCHAR(255) NOT NULL,
    enrolled_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    course_progress_percentage DECIMAL(5, 2) DEFAULT 0.0,
    course_progress_status VARCHAR(50) NOT NULL DEFAULT 'enrolled',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, course_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE INDEX idx_enrollments_user_id ON enrollments(user_id);
CREATE INDEX idx_enrollments_course_id ON enrollments(course_id);
CREATE INDEX idx_enrollments_user_course ON enrollments(user_id, course_id);

-- Module progress table
CREATE TABLE IF NOT EXISTS module_progress (
    enrollment_id VARCHAR(255) NOT NULL,
    module_id VARCHAR(255) NOT NULL,
    progress_percentage DECIMAL(5, 2) DEFAULT 0.0,
    progress_status VARCHAR(50) NOT NULL DEFAULT 'started',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (enrollment_id, module_id),
    FOREIGN KEY (enrollment_id) REFERENCES enrollments(id) ON DELETE CASCADE,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE
);

CREATE INDEX idx_module_progress_enrollment ON module_progress(enrollment_id);

-- Lesson progress table
CREATE TABLE IF NOT EXISTS lesson_progress (
    enrollment_id VARCHAR(255) NOT NULL,
    lesson_id VARCHAR(255) NOT NULL,
    progress_percentage DECIMAL(5, 2) DEFAULT 0.0,
    progress_status VARCHAR(50) NOT NULL DEFAULT 'started',
    exercise_score DECIMAL(5, 2) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (enrollment_id, lesson_id),
    FOREIGN KEY (enrollment_id) REFERENCES enrollments(id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);

CREATE INDEX idx_lesson_progress_enrollment ON lesson_progress(enrollment_id);
