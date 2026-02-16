-- +goose Up
-- +goose StatementBegin
CREATE TABLE visits (
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT REFERENCES links(id) NOT NULL,
    ip VARCHAR(45),
    user_agent VARCHAR(255),
    referer VARCHAR(2083),
    "status" SMALLINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_visits_link_id ON visits(link_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS visits;
-- +goose StatementEnd
