-- +migrate UP
CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    login TEXT NOT NULL,
    hash_password TEXT NOT NULL,
    email TEXT
);