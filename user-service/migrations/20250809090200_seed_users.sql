-- +goose Up
INSERT INTO users (username, email, password_hash)
VALUES ('demo', 'demo@example.com', '$2y$12$PLACEHOLDER_BCRYPT_HASH');
-- +goose Down
DELETE FROM users WHERE email = 'demo@example.com';
