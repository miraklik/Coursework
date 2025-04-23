CREATE TABLE books (
    ID BIGSERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author_id BIGINT NOT NULL,
    FOREIGN KEY (author_id) REFERENCES authors(ID)
);

CREATE INDEX idx_books_title ON books(title);
