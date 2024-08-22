package repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	commonDs "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/data_source"
	common "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	"github.com/NidzamuddinMuzakki/movies-abishar/model"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	CreateUsers(ctx context.Context, tx *sqlx.Tx, payload model.UsersModel) (id int64, err error)
	LoginUsers(ctx context.Context, payload model.UsersModel) (result *model.UsersModel, err error)
}

type users struct {
	common common.IRegistry
	master *sqlx.DB
	slave  *sqlx.DB
}

func NewUsersRepository(common common.IRegistry, master *sqlx.DB, slave *sqlx.DB) IUsersRepository {
	return &users{
		common: common,
		master: master,
		slave:  slave,
	}
}

func (r users) CreateUsers(ctx context.Context, tx *sqlx.Tx, data model.UsersModel) (id int64, err error) {
	password := []byte(data.Password)
	passwordMD5 := fmt.Sprintf("%x", md5.Sum(password))
	insertQuery := "insert into users (username,password, created_at, updated_at) VALUES(?,?,?,?)"
	stmtx, err := tx.PreparexContext(ctx, insertQuery)
	if err != nil {

		return 0, err
	}
	res, err := stmtx.ExecContext(
		ctx,
		data.Username,
		passwordMD5,
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

func (r users) LoginUsers(ctx context.Context, data model.UsersModel) (result *model.UsersModel, err error) {
	var (
		it model.UsersModel
	)
	password := []byte(data.Password)
	passwordMD5 := fmt.Sprintf("%x", md5.Sum(password))

	selectQuery := "select username, password from users where username = ? and password= ? "

	err = commonDs.Exec(ctx, r.master, commonDs.NewStatement(&it, selectQuery, data.Username, passwordMD5))
	if err != nil {
		return nil, err
	}

	return &it, nil

}
