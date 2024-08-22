package service

import (
	"github.com/NidzamuddinMuzakki/movies-abishar/service/health"
)

// @Notice: Register your services here

type IRegistry interface {
	GetHealth() health.IHealth
	GetUsersService() IUsersService
	GetMoviesService() IMoviesService
}

type Registry struct {
	health        health.IHealth
	usersService  IUsersService
	moviesService IMoviesService
}

func NewRegistry(health health.IHealth, usersService IUsersService, moviesService IMoviesService) *Registry {
	return &Registry{
		health:        health,
		usersService:  usersService,
		moviesService: moviesService,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}

func (r *Registry) GetUsersService() IUsersService {
	return r.usersService
}

func (r *Registry) GetMoviesService() IMoviesService {
	return r.moviesService
}
