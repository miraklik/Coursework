CREATE TABLE authors (
    ID BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE INDEX idx_authors_name ON authors(name);