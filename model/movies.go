package model

import (
	"mime/multipart"
	"time"
)

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

type MoviesGenre struct {
	Genre string `json:"genre" db:"genre"`
	Count int    `json:"count" db:"count"`
}

type MoviesView struct {
	Id    string `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Count int    `json:"count" db:"count"`
}
type RequestMoviesModel struct {
	Title       string                `form:"title" db:"title" validate:"required"`
	Description string                `form:"description" db:"description" validate:"required"`
	Duration    int                   `form:"duration" db:"duration" validate:"required,number,min=1"`
	Artists     string                `form:"artists" db:"artists" validate:"required"`
	Genres      string                `form:"genres" db:"genres" validate:"required"`
	UrlWatch    *multipart.FileHeader `form:"url_watch" db:"url_watch" validate:"required"`
}

type RequestUpdateMoviesModel struct {
	Id          int                   `uri:"id" db:"id" validate:"required,min=1"`
	Title       string                `form:"title" db:"title" validate:"required"`
	Description string                `form:"description" db:"description" validate:"required"`
	Duration    int                   `form:"duration" db:"duration" validate:"required,number,min=1"`
	Artists     string                `form:"artists" db:"artists" validate:"required"`
	Genres      string                `form:"genres" db:"genres" validate:"required"`
	UrlWatch    *multipart.FileHeader `form:"url_watch" db:"url_watch"`
}

type RequestGetListMoviesModel struct {
	Search string `form:"search"`
	Limit  uint   `form:"limit" validate:"required,number,min=1"`
	Offset uint   `form:"offset" validate:"required,number,min=1"`
}

type RequestGetDetailMoviesModel struct {
	Id uint `uri:"id" validate:"required,number,min=1"`
}

type RequestGetDetailViewMoviesModel struct {
	Id       uint `uri:"id" validate:"required,number,min=1"`
	Duration uint `uri:"duration" validate:"required,number,min=1"`
}
