package service

import (
	"context"
	"fmt"
	"time"

	"github.com/NidzamuddinMuzakki/movies-abishar/common/util"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/logger"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/movies-abishar/model"
	"github.com/NidzamuddinMuzakki/movies-abishar/repository"
	"github.com/jmoiron/sqlx"
)

type IUsersService interface {
	CreateUsers(ctx context.Context, payload model.RequestUsersModel) error
	LoginUsers(ctx context.Context, payload model.RequestUsersModel) error
}

type usersService struct {
	common       registry.IRegistry
	repoRegistry repository.IRegistry
}

func NewUsersService(common registry.IRegistry, repoRegistry repository.IRegistry) IUsersService {
	return &usersService{
		common:       common,
		repoRegistry: repoRegistry,
	}
}

func (s usersService) CreateUsers(ctx context.Context, payload model.RequestUsersModel) error {
	now := time.Now()

	doFunc := util.TxFunc(func(tx *sqlx.Tx) error {
		data := model.UsersModel{
			Username:  payload.Username,
			Password:  payload.Password,
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		_, err := s.repoRegistry.GetUsersRepository().CreateUsers(ctx, tx, data)
		fmt.Println(err, "err2")
		if err != nil {
			logger.Error(ctx, err.Error(), err)
			return err
		}
		return nil
	})
	err := s.repoRegistry.GetUtilTx().WithTx(ctx, doFunc, nil)
	fmt.Println(err, "err")
	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	return nil
}

func (s usersService) LoginUsers(ctx context.Context, payload model.RequestUsersModel) error {
	data := model.UsersModel{
		Username: payload.Username,
		Password: payload.Password,
	}
	_, err := s.repoRegistry.GetUsersRepository().LoginUsers(ctx, data)

	if err != nil {
		logger.Error(ctx, err.Error(), err)
		return err
	}
	return nil
}
