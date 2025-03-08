-- +goose Up
ALTER TABLE notes
ADD CONSTRAINT unique_user_note UNIQUE (user_id, body);
-- +goose Down
ALTER TABLE notes DROP CONSTRAINT unique_user_note;