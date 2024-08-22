package health

import (
	"context"

	"github.com/NidzamuddinMuzakki/movies-abishar/model"

	"github.com/jmoiron/sqlx"
)

const (
	OK  = "OK"
	BAD = "BAD"
)

type IHealth interface {
	Check(ctx context.Context) model.HTTPResponse
}

type Health struct {
	master *sqlx.DB
	slave  *sqlx.DB
}

func NewHealth(master, slave *sqlx.DB) *Health {
	return &Health{
		master: master,
		slave:  slave,
	}
}

func (s *Health) Check(ctx context.Context) model.HTTPResponse {
	var response = model.HTTPResponse{
		Master: OK,
		Slave:  OK,
	}

	return response
}
