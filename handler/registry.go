package handler

import (
	"github.com/NidzamuddinMuzakki/movies-abishar/handler/health"
)

// @Notice: Register your http deliveries here

type IRegistry interface {
	GetHealth() health.IHealth
	GetUsers() IUsers
}

type Registry struct {
	health health.IHealth
	users  IUsers
}

func NewRegistry(health health.IHealth, users IUsers) *Registry {
	return &Registry{
		health: health,
		users:  users,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}

func (r *Registry) GetUsers() IUsers {
	return r.users
}
