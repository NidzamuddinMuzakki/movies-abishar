package repository

import "github.com/NidzamuddinMuzakki/movies-abishar/common/util"

// @Notice: Register your repositories here

type IRegistry interface {
	GetUsersRepository() IUsersRepository
	GetUtilTx() *util.TransactionRunner
}

type Registry struct {
	usersRepository IUsersRepository
	masterUtilTx    *util.TransactionRunner
}

func NewRegistryRepository(
	masterUtilTx *util.TransactionRunner,
	usersRepository IUsersRepository,
) *Registry {
	return &Registry{
		masterUtilTx:    masterUtilTx,
		usersRepository: usersRepository,
	}
}

func (r Registry) GetUtilTx() *util.TransactionRunner {
	return r.masterUtilTx
}

func (r Registry) GetUsersRepository() IUsersRepository {
	return r.usersRepository
}
