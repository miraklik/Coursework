CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    group_name VARCHAR(50) NOT NULL
);

CREATE INDEX idx_students_group ON students(group_name);

CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(200)
);

CREATE INDEX idx_organizations_name ON organizations(name);