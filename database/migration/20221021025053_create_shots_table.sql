-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shots (
  id serial,
  title VARCHAR(255),
  img VARCHAR(255),
  description VARCHAR(255),
  category VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shots;
-- +goose StatementEnd
