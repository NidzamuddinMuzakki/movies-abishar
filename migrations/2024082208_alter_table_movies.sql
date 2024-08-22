-- +goose Up
-- +goose StatementBegin
-- abishar.count_viewed_genre definition
ALTER TABLE movies
ADD COLUMN url_watch TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
DROP COLUMN url_watch;
-- +goose StatementEnd
