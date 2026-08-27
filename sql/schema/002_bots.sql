-- +goose Up
CREATE TABLE bots (
    id UUID PRIMARY KEY,
    bot_id int NOT NULL,
    skills int[] NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

-- +goose Down
DROP TABLE bots;