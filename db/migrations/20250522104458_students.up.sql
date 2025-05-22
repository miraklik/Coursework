CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    group_name VARCHAR(50) NOT NULL
);

CREATE INDEX idx_students_group ON students(group_name);