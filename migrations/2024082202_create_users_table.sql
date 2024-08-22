-- +goose Up
-- +goose StatementBegin
-- abishar.users definition
CREATE TABLE IF NOT EXISTS `users`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT ,
    `username` varchar(255),
    `password` varchar(255),
    `created_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT UNIQUE_USERS UNIQUE (username)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `users`;
-- +goose StatementEnd
