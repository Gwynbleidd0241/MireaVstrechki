-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    creator_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
