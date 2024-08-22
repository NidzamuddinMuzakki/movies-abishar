package service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NidzamuddinMuzakki/movies-abishar/common/util"
	"github.com/NidzamuddinMuzakki/movies-abishar/config"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/logger"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/movies-abishar/model"
	"github.com/NidzamuddinMuzakki/movies-abishar/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func checkDuplicateInt(intSlice []string) bool {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	if len(intSlice) == len(list) {
		return false
	}
	return true
}

type IMoviesService interface {
	CreateMovies(ctx context.Context, payload model.RequestMoviesModel) error
	UpdateMovies(ctx context.Context, payload model.RequestUpdateMoviesModel) error
	GetMoviesList(ctx context.Context, payload model.RequestGetListMoviesModel) ([]model.MoviesModel, uint64, error)
	GetMoviesById(ctx context.Context, id int) (*model.MoviesModel, error)

	GetMoviesByIdUserView(ctx context.Context, id int, username string, duration int) (*model.MoviesModel, error)
	MoviesVote(ctx context.Context, id int, username string) error
	UnVoteMovies(ctx context.Context, id int, username string) error

	GetVoteMovies(ctx context.Context, username string) ([]model.MoviesModel, error)

	GetMostVoteMovies(ctx context.Context) (*model.MoviesView, error)
	GetMostViewedMovies(ctx context.Context) (*model.MoviesView, error)
	GetMostViewedGenre(ctx context.Context) (*model.MoviesGenre, error)
}

type moviesService struct {
	common       registry.IRegistry
	repoRegistry repository.IRegistry
}

func NewMoviesService(common registry.IRegistry, repoRegistry repository.IRegistry) IMoviesService {
	return &moviesService{
		common:       common,
		repoRegistry: repoRegistry,
	}
}

func (s moviesService) GetVoteMovies(ctx context.Context, username string) ([]model.MoviesModel, error) {
	dataUser, err := s.repoRegistry.GetUsersRepository().GetUsersByUsername(ctx, username)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}
	resp, err := s.repoRegistry.GetMoviesRepository().GetVoteMovies(ctx, dataUser.Id)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	return resp, nil

}
func (s moviesService) GetMostVoteMovies(ctx context.Context) (*model.MoviesView, error) {

	resp, err := s.repoRegistry.GetMoviesRepository().GetMostVoteMovies(ctx)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	return &resp, nil

}
func (s moviesService) GetMostViewedMovies(ctx context.Context) (*model.MoviesView, error) {

	resp, err := s.repoRegistry.GetMoviesRepository().GetMostViewedMovies(ctx)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	return &resp, nil

}
func (s moviesService) GetMostViewedGenre(ctx context.Context) (*model.MoviesGenre, error) {

	resp, err := s.repoRegistry.GetMoviesRepository().GetMostViewedGenre(ctx)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	return &resp, nil

}
func (s moviesService) GetMoviesById(ctx context.Context, id int) (*model.MoviesModel, error) {

	resp, err := s.repoRegistry.GetMoviesRepository().GetMoviesById(ctx, id)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	return &resp, nil

}

func (s moviesService) MoviesVote(ctx context.Context, id int, username string) error {

	resp, err := s.repoRegistry.GetMoviesRepository().GetMoviesById(ctx, id)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}

	idUser, err := s.repoRegistry.GetUsersRepository().GetUsersByUsername(ctx, username)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	doFunc := util.TxFunc(func(tx *sqlx.Tx) error {
		_, err := s.repoRegistry.GetMoviesRepository().InsertVoteMovie(ctx, tx, resp.Id, idUser.Id)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}

		_, err = s.repoRegistry.GetMoviesRepository().UpsertVoteCountMovie(ctx, tx, resp.Id)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}

		return nil
	})
	err = s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)

	if err != nil {

		logger.Error(ctx, err.Error(), err)
		return err
	}

	return nil

}

func (s moviesService) UnVoteMovies(ctx context.Context, id int, username string) error {

	resp, err := s.repoRegistry.GetMoviesRepository().GetMoviesById(ctx, id)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}

	idUser, err := s.repoRegistry.GetUsersRepository().GetUsersByUsername(ctx, username)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	doFunc := util.TxFunc(func(tx *sqlx.Tx) error {
		_, err := s.repoRegistry.GetMoviesRepository().UpdateVoteMovie(ctx, tx, resp.Id, idUser.Id)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}

		_, err = s.repoRegistry.GetMoviesRepository().UpdateVoteCountMovie(ctx, tx, resp.Id)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}

		return nil
	})
	err = s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)

	if err != nil {

		logger.Error(ctx, err.Error(), err)
		return err
	}

	return nil

}

func (s moviesService) GetMoviesByIdUserView(ctx context.Context, id int, username string, duration int) (*model.MoviesModel, error) {

	resp, err := s.repoRegistry.GetMoviesRepository().GetMoviesById(ctx, id)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	idUser, err := s.repoRegistry.GetUsersRepository().GetUsersByUsername(ctx, username)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	genres := strings.Split(resp.Genres, ",")
	doFunc := util.TxFunc(func(tx *sqlx.Tx) error {
		_, err := s.repoRegistry.GetMoviesRepository().UpsertViewedMovie(ctx, tx, resp.Id, idUser.Id, duration)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}

		_, err = s.repoRegistry.GetMoviesRepository().UpsertViewedCountMovie(ctx, tx, resp.Id)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}
		for _, genre := range genres {
			_, err = s.repoRegistry.GetMoviesRepository().UpsertViewedCountGenres(ctx, tx, genre)
			if err != nil {
				logger.Error(ctx, err.Error(), err)
				return err
			}
		}
		return nil
	})
	err = s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)

	if err != nil {

		logger.Error(ctx, err.Error(), err)
		return nil, err
	}

	return &resp, nil

}

func (s moviesService) UpdateMovies(ctx context.Context, payload model.RequestUpdateMoviesModel) error {
	now := time.Now()
	uuids := uuid.New()
	artist := strings.Split(payload.Artists, ",")
	genres := strings.Split(payload.Genres, ",")
	duplicateArtis := checkDuplicateInt(artist)
	if duplicateArtis {
		err := errors.New("duplicateArtis")
		logger.Error(ctx, err.Error(), err)
		return err
	}
	duplicateGenres := checkDuplicateInt(genres)
	if duplicateGenres {
		err := errors.New("duplicateGenres")
		logger.Error(ctx, err.Error(), err)
		return err
	}
	var tempFile *os.File
	var da multipart.File
	if payload.UrlWatch != nil {
		FILEnAME := (uuids.String() + "*" + filepath.Ext(payload.UrlWatch.Filename))
		// FILEnAME2 := (uuids.String() + filepath.Ext(file.Filename))
		tempFile, err := ioutil.TempFile("files", FILEnAME)
		defer tempFile.Close()
		da, err = payload.UrlWatch.Open()
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}

	}
	getData, err := s.repoRegistry.GetMoviesRepository().GetMoviesById(ctx, payload.Id)
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	doFunc := util.TxFunc(func(tx *sqlx.Tx) error {
		url := getData.UrlWatch
		if payload.UrlWatch != nil {
			url = config.Cold.AppHost + "/files/" + strings.Split(tempFile.Name(), "\\")[1]

		}
		data := model.MoviesModel{
			Id:          payload.Id,
			Title:       payload.Title,
			Description: payload.Description,
			Duration:    payload.Duration,
			Artists:     payload.Artists,
			Genres:      payload.Genres,
			UrlWatch:    url,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}
		_, err := s.repoRegistry.GetMoviesRepository().UpdateMovies(ctx, tx, data)

		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}
		if payload.UrlWatch != nil {
			fileBytes, err := ioutil.ReadAll(da)
			if err != nil {
				logger.Error(ctx, err.Error(), err)
				return err
			}
			_, err = tempFile.Write(fileBytes)
			if err != nil {
				logger.Error(ctx, err.Error(), err)
				return err
			}

		}

		return nil
	})
	err = s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)
	fmt.Println(err, "err")
	if err != nil {
		if payload.UrlWatch != nil {
			tempFile.Close()
			errs := os.Remove(tempFile.Name())
			fmt.Println(errs, "err remove")
		}

		logger.Error(ctx, err.Error(), err)
		return err
	}
	return nil
}
func (s moviesService) CreateMovies(ctx context.Context, payload model.RequestMoviesModel) error {
	now := time.Now()
	uuids := uuid.New()
	artist := strings.Split(payload.Artists, ",")
	genres := strings.Split(payload.Genres, ",")
	duplicateArtis := checkDuplicateInt(artist)
	if duplicateArtis {
		err := errors.New("duplicateArtis")
		logger.Error(ctx, err.Error(), err)
		return err
	}
	duplicateGenres := checkDuplicateInt(genres)
	if duplicateGenres {
		err := errors.New("duplicateGenres")
		logger.Error(ctx, err.Error(), err)
		return err
	}
	FILEnAME := (uuids.String() + "*" + filepath.Ext(payload.UrlWatch.Filename))
	// FILEnAME2 := (uuids.String() + filepath.Ext(file.Filename))
	tempFile, err := ioutil.TempFile("files", FILEnAME)
	defer tempFile.Close()
	da, err := payload.UrlWatch.Open()
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}

	doFunc := util.TxFunc(func(tx *sqlx.Tx) error {
		data := model.MoviesModel{
			Title:       payload.Title,
			Description: payload.Description,
			Duration:    payload.Duration,
			Artists:     payload.Artists,
			Genres:      payload.Genres,
			UrlWatch:    config.Cold.AppHost + "/files/" + strings.Split(tempFile.Name(), "\\")[1],
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}
		_, err := s.repoRegistry.GetMoviesRepository().CreateMovies(ctx, tx, data)
		fmt.Println(err, "err2")
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}

		fileBytes, err := ioutil.ReadAll(da)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}
		_, err = tempFile.Write(fileBytes)
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}
		return nil
	})
	err = s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)
	fmt.Println(err, "err")
	if err != nil {
		tempFile.Close()
		errs := os.Remove(tempFile.Name())
		fmt.Println(errs, "err remove")
		logger.Error(ctx, err.Error(), err)
		return err
	}
	return nil
}

func (s moviesService) GetMoviesList(ctx context.Context, payload model.RequestGetListMoviesModel) ([]model.MoviesModel, uint64, error) {

	if payload.Search != "" {
		payload.Search = strings.ToLower(payload.Search)
	}

	resp, totalCounts, err := s.repoRegistry.GetMoviesRepository().GetListMovies(ctx, payload)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return nil, 0, err
	}

	return resp, totalCounts, nil

}
