-- +goose Up
-- +goose StatementBegin
-- abishar.movies definition
CREATE TABLE IF NOT EXISTS `movies`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT ,
    `title`       varchar(200)  NOT NULL DEFAULT '',
    `description`       varchar(200)  NOT NULL DEFAULT '',
    `duration` int  NOT NULL DEFAULT 0,
    `artists` TEXT NOT NULL DEFAULT '',
    `genres` TEXT NOT NULL DEFAULT '',
    `created_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT UNIQUE_MOVIES UNIQUE (title)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `movies`;
-- +goose StatementEnd
