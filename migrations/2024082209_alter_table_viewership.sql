-- +goose Up
-- +goose StatementBegin
-- abishar.count_viewed_genre definition
ALTER TABLE viewership
ADD COLUMN duration int NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE viewership
DROP COLUMN duration;
-- +goose StatementEnd
