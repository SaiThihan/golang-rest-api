-- +goose Up
-- +goose StatementBegin
CREATE TABLE tokens (
    Hash BYTEA PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens
-- +goose StatementEnd