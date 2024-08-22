-- +goose Up
-- +goose StatementBegin
-- abishar.vote_movies definition
CREATE TABLE IF NOT EXISTS `vote_movies`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT ,
    `movie_id` int NOT NULL DEFAULT 0,
    `user_id` int NOT NULL DEFAULT 0,
    `created_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT UNIQUE_VOTES_MOVIES UNIQUE (movie_id,user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `vote_movies`;
-- +goose StatementEnd
