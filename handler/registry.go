package handler

import (
	"github.com/NidzamuddinMuzakki/movies-abishar/handler/health"
)

// @Notice: Register your http deliveries here

type IRegistry interface {
	GetHealth() health.IHealth
	GetUsers() IUsers
	GetMovies() IMovies
}

type Registry struct {
	health health.IHealth
	users  IUsers
	movies IMovies
}

func NewRegistry(health health.IHealth, users IUsers, movies IMovies) *Registry {
	return &Registry{
		health: health,
		users:  users,
		movies: movies,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}

func (r *Registry) GetUsers() IUsers {
	return r.users
}

func (r *Registry) GetMovies() IMovies {
	return r.movies
}
