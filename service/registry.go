package service

import (
	"github.com/NidzamuddinMuzakki/movies-abishar/service/health"
)

// @Notice: Register your services here

type IRegistry interface {
	GetHealth() health.IHealth
	GetUsersService() IUsersService
}

type Registry struct {
	health       health.IHealth
	usersService IUsersService
}

func NewRegistry(health health.IHealth, usersService IUsersService) *Registry {
	return &Registry{
		health:       health,
		usersService: usersService,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}

func (r *Registry) GetUsersService() IUsersService {
	return r.usersService
}
