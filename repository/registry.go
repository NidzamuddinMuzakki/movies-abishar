package repository

import "github.com/NidzamuddinMuzakki/movies-abishar/common/util"

// @Notice: Register your repositories here

type IRegistry interface {
	GetUsersRepository() IUsersRepository
	GetMoviesRepository() IMoviesRepository
	GetUtilTx() *util.TransactionRunner
}

type Registry struct {
	moviesRepository IMoviesRepository
	usersRepository  IUsersRepository
	masterUtilTx     *util.TransactionRunner
}

func NewRegistryRepository(
	masterUtilTx *util.TransactionRunner,
	usersRepository IUsersRepository,
	moviesRepository IMoviesRepository,
) *Registry {
	return &Registry{
		masterUtilTx:     masterUtilTx,
		usersRepository:  usersRepository,
		moviesRepository: moviesRepository,
	}
}

func (r Registry) GetUtilTx() *util.TransactionRunner {
	return r.masterUtilTx
}

func (r Registry) GetUsersRepository() IUsersRepository {
	return r.usersRepository
}

func (r Registry) GetMoviesRepository() IMoviesRepository {
	return r.moviesRepository
}
