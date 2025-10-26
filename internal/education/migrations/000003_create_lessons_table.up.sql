CREATE TABLE IF NOT EXISTS lessons (
    id VARCHAR(255) PRIMARY KEY,
    module_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    overview TEXT,
    content TEXT,
    video_id VARCHAR(255),
    duration INTEGER NOT NULL DEFAULT 0,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_lessons_module_id ON lessons(module_id);
CREATE INDEX IF NOT EXISTS idx_lessons_order ON lessons(module_id, order_index);
