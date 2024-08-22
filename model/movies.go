package model

import "time"

type MoviesModel struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Duration    int    `json:"duration" db:"duration"`
	Artists     string `json:"artists" db:"artists"`
	Genres      string `json:"genres" db:"genres"`
	UrlWatch    string `json:"url_watch" db:"url_watch"`

	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type RequestMoviesModel struct {
	Username string `json:"username" db:"username" validate:"required,alphanum"`
	Password string `json:"password" db:"password" validate:"required,alphanum"`
}
