-- +goose Up
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE CHECK (email LIKE '%@gmail.com'),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    password TEXT NOT NULL CHECK (LENGTH(password) > 8),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);
-- Create notes table
CREATE TABLE IF NOT EXISTS notes (
    notes_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id INT NOT NULL,
    body VARCHAR NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_note UNIQUE (user_id, body)
);
-- Index for performance
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
-- +goose Down
-- Drop tables if rolling back
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS users;