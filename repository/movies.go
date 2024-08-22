package repository

import (
	"context"
	"errors"

	common "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/movies-abishar/model"
	"github.com/jmoiron/sqlx"
)

type IMoviesRepository interface {
	CreateMovies(ctx context.Context, tx *sqlx.Tx, payload model.MoviesModel) (id int64, err error)
}

type movies struct {
	common common.IRegistry
	master *sqlx.DB
	slave  *sqlx.DB
}

func NewMoviesRepository(common common.IRegistry, master *sqlx.DB, slave *sqlx.DB) IMoviesRepository {
	return &movies{
		common: common,
		master: master,
		slave:  slave,
	}
}
func (r movies) CreateMovies(ctx context.Context, tx *sqlx.Tx, data model.MoviesModel) (id int64, err error) {
	insertQuery := "insert into movies (title,description,duration,artists,genres,url_watch, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?)"
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		data.Title,
		data.Description,
		data.Duration,
		data.Artists,
		data.Genres,
		data.UrlWatch,
		data.CreatedAt,
		data.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}
