--@block
-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     email VARCHAR(255) NOT NULL UNIQUE,
--     password VARCHAR(256) NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );
-- CREATE TABLE notes (
--     id SERIAL PRIMARY KEY,
--     body VARCHAR(255) NOT NULL,
--     deleted_at TIMESTAMP DEFAULT NULL,
--     user_id INT NOT NULL,
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
-- );
-- CREATE INDEX idx_notes_user_id ON notes(user_id);
--@block
SELECT *
FROM notes;