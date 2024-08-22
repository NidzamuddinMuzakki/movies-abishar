package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	commonDs "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/data_source"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/logger"
	common "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/movies-abishar/model"
	"github.com/jmoiron/sqlx"
)

type IMoviesRepository interface {
	CreateMovies(ctx context.Context, tx *sqlx.Tx, payload model.MoviesModel) (id int64, err error)
	UpdateMovies(ctx context.Context, tx *sqlx.Tx, payload model.MoviesModel) (id int64, err error)
	GetListMovies(ctx context.Context, payload model.RequestGetListMoviesModel) ([]model.MoviesModel, uint64, error)
	GetMoviesById(ctx context.Context, id int) (model.MoviesModel, error)
	UpsertViewedMovie(ctx context.Context, tx *sqlx.Tx, movie_id int, user_id int, duration int) (id int64, err error)
	UpsertViewedCountMovie(ctx context.Context, tx *sqlx.Tx, movie_id int) (id int64, err error)
	UpsertViewedCountGenres(ctx context.Context, tx *sqlx.Tx, genre string) (id int64, err error)

	InsertVoteMovie(ctx context.Context, tx *sqlx.Tx, movie_id int, user_id int) (id int64, err error)
	UpsertVoteCountMovie(ctx context.Context, tx *sqlx.Tx, movie_id int) (id int64, err error)

	UpdateVoteMovie(ctx context.Context, tx *sqlx.Tx, movie_id int, user_id int) (id int64, err error)
	UpdateVoteCountMovie(ctx context.Context, tx *sqlx.Tx, movie_id int) (id int64, err error)

	GetVoteMovies(ctx context.Context, user_id int) ([]model.MoviesModel, error)
	GetMostViewedMovies(ctx context.Context) (model.MoviesView, error)
	GetMostVoteMovies(ctx context.Context) (model.MoviesView, error)
	GetMostViewedGenre(ctx context.Context) (model.MoviesGenre, error)
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

func (r movies) GetVoteMovies(ctx context.Context, user_id int) ([]model.MoviesModel, error) {
	var (
		data []model.MoviesModel
	)

	selectQuery := `select 
    	m.id,
		m.title, 
		m.description, 
		m.duration, 
		m.artists,
		m.genres,
		m.url_watch,
		m.created_at,
		m.updated_at
	from movies m INNER JOIN vote_movies vtm ON m.id=vtm.movie_id where vtm.user_id = ? `

	err := commonDs.Exec(ctx, r.master, commonDs.NewStatement(&data, selectQuery, user_id))

	if err != nil {
		return data, err
	}

	return data, nil
}
func (r movies) GetMostViewedMovies(ctx context.Context) (model.MoviesView, error) {
	var (
		data model.MoviesView
	)

	selectQuery := `select 
    	ms.id,
		ms.title,
		cvw.count
	from movies ms INNER JOIN count_viewership cvw ON ms.id=cvw.movie_id order by cvw.count desc,ms.id desc limit 1`

	err := commonDs.Exec(ctx, r.master, commonDs.NewStatement(&data, selectQuery))
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r movies) GetMostVoteMovies(ctx context.Context) (model.MoviesView, error) {
	var (
		data model.MoviesView
	)

	selectQuery := `select 
    	ms.id,
		ms.title,
		cvw.count
	from movies ms INNER JOIN count_vote_movies cvw ON ms.id=cvw.movie_id order by cvw.count desc,ms.id desc limit 1`

	err := commonDs.Exec(ctx, r.master, commonDs.NewStatement(&data, selectQuery))
	if err != nil {
		return data, err
	}

	return data, nil
}
func (r movies) GetMostViewedGenre(ctx context.Context) (model.MoviesGenre, error) {
	var (
		data model.MoviesGenre
	)

	selectQuery := `select 
    	genre,
		count
	from count_viewed_genre  order by count desc, id desc limit 1`

	err := commonDs.Exec(ctx, r.master, commonDs.NewStatement(&data, selectQuery))
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r movies) UpsertVoteCountMovie(ctx context.Context, tx *sqlx.Tx, movie_id int) (id int64, err error) {
	now := time.Now()
	insertQuery := "INSERT INTO count_vote_movies (movie_id, count,created_at,updated_at) " +
		"VALUES (?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE " +
		"count = count+1 ," +
		"updated_at = ? "
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		movie_id,
		1,
		now,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}
func (r movies) InsertVoteMovie(ctx context.Context, tx *sqlx.Tx, movie_id int, user_id int) (id int64, err error) {
	now := time.Now()
	insertQuery := "INSERT INTO vote_movies (movie_id, user_id,created_at,updated_at) " +
		"VALUES (?,?,?,?) "
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		movie_id,
		user_id,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) UpdateVoteCountMovie(ctx context.Context, tx *sqlx.Tx, movie_id int) (id int64, err error) {
	now := time.Now()
	insertQuery := "UPDATE  count_vote_movies set count=count-1,updated_at=? where movie_id=? and count-1>-1"
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		now,
		movie_id,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	fmt.Println(intss, "update count")
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) UpdateVoteMovie(ctx context.Context, tx *sqlx.Tx, movie_id int, user_id int) (id int64, err error) {
	insertQuery := "delete from vote_movies where movie_id=? and user_id=? "
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		movie_id,
		user_id,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	fmt.Println(intss, "update delete")

	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) UpsertViewedMovie(ctx context.Context, tx *sqlx.Tx, movie_id int, user_id int, duration int) (id int64, err error) {
	now := time.Now()
	insertQuery := "INSERT INTO viewership (movie_id, user_id,duration,created_at,updated_at) " +
		"VALUES (?,?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE " +
		"duration = ? ," +
		"updated_at = ? "
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		movie_id,
		user_id,
		duration,
		now,
		now,
		duration,
		now,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) UpsertViewedCountMovie(ctx context.Context, tx *sqlx.Tx, movie_id int) (id int64, err error) {
	now := time.Now()
	insertQuery := "INSERT INTO count_viewership (movie_id, count,created_at,updated_at) " +
		"VALUES (?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE " +
		"count = count+1 ," +
		"updated_at = ? "
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		movie_id,
		1,
		now,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) UpsertViewedCountGenres(ctx context.Context, tx *sqlx.Tx, genre string) (id int64, err error) {
	now := time.Now()
	insertQuery := "INSERT INTO count_viewed_genre (genre, count,created_at,updated_at) " +
		"VALUES (?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE " +
		"count = count+1, " +
		"updated_at = ? "
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		genre,
		1,
		now,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) GetMoviesById(ctx context.Context, id int) (model.MoviesModel, error) {
	var (
		data model.MoviesModel
	)

	selectQuery := `select 
    	id,
		title, 
		description, 
		duration, 
		artists,
		genres,
		url_watch,
		created_at,
		updated_at
	from movies
	where id = ? `

	err := commonDs.Exec(ctx, r.master, commonDs.NewStatement(&data, selectQuery, id))

	if err != nil {
		return data, err
	}

	return data, nil
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
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) UpdateMovies(ctx context.Context, tx *sqlx.Tx, data model.MoviesModel) (id int64, err error) {
	insertQuery := "update movies set title=?,description=?,duration=?,artists=?,genres=?,url_watch=?,updated_at=? where id=?"
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	fmt.Println(err, data.Id, "err 2")
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
		data.UpdatedAt,
		data.Id,
	)
	fmt.Println(err, "err 2")
	if err != nil {
		return 0, err
	}
	intss, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if intss == 0 {
		err = errors.New("nothing update")
		return 0, err
	}

	return intss, nil
}

func (r movies) GetListMovies(ctx context.Context, payload model.RequestGetListMoviesModel) ([]model.MoviesModel, uint64, error) {

	var (
		list              []model.MoviesModel
		totalTransactions uint64
		filters           []string
		args              []any
	)
	countQuery := "select count(id) from movies where 1=1 "
	selectQuery := "SELECT " +
		"id, " +
		"title, " +
		"description, " +
		"duration, " +
		"artists, " +
		"genres, " +
		"url_watch, " +

		"created_at, " +
		"updated_at " +
		"FROM movies " +
		"where 1=1 "
	if payload.Search != "" {
		filters = append(filters, "and  (LOWER(title) like ? or LOWER(description) like ? or LOWER(artists) like ? or LOWER(genres) like ? )  ")
		args = append(args, "%"+payload.Search+"%", "%"+payload.Search+"%", "%"+payload.Search+"%", "%"+payload.Search+"%")
	}

	for _, f := range filters {
		countQuery = fmt.Sprintf("%s %s", countQuery, f)
		selectQuery = fmt.Sprintf("%s %s", selectQuery, f)
	}

	offset := (payload.Limit * payload.Offset) - payload.Limit
	selectQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", selectQuery, payload.Limit, offset)

	err := commonDs.Exec(ctx, r.slave, commonDs.NewStatement(&list, selectQuery, args...))
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, 0, err
	}

	err = commonDs.Exec(ctx, r.slave, commonDs.NewStatement(&totalTransactions, countQuery, args...))
	if err != nil {

		logger.Error(ctx, err.Error(), err)
		return nil, 0, err
	}

	return list, totalTransactions, nil

}
