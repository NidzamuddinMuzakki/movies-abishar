-- +goose Up
-- +goose StatementBegin
-- abishar.count_viewed_genre definition
CREATE TABLE IF NOT EXISTS `count_viewed_genre`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT ,
    `genre` varchar(255) NOT NULL DEFAULT '',
    `count` int NOT NULL DEFAULT 0,
    `created_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT UNIQUE_COUNT_VIEWED_GENRE UNIQUE (genre)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `count_viewed_genre`;
-- +goose StatementEnd
