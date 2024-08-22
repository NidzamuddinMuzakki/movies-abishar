-- +goose Up
-- +goose StatementBegin
-- abishar.count_viewership definition
CREATE TABLE IF NOT EXISTS `count_viewership`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT ,
    `movie_id` int NOT NULL DEFAULT 0,
    `count` int NOT NULL DEFAULT 0,
    `created_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT UNIQUE_COUNT_VIEWERSHIP UNIQUE (movie_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `count_viewership`;
-- +goose StatementEnd
