-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    original_url VARCHAR(2083) NOT NULL,
    short_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_links_short_name ON links(short_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd
