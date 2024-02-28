-- +goose Up
CREATE TABLE logged_out (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID UNIQUE NOT NULL,
    tokens BYTEA NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
);

-- +goose Down
DROP TABLE logged_out;