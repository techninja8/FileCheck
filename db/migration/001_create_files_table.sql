CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    hash VARCHAR(64) NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    location VARCHAR(1024) NOT NULL
);