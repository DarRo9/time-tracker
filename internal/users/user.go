package users

import (
	"context"
)

type userId struct {
	Id int `json:"id"`
}

type passport struct {
	PassportNumber string `json:"passportNumber"`
}

type user struct {
	Id         int    `gorm:"column:id"`
	Passport   string `json:"passport"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

type searchReq struct {
	page     int
	pageSize int
	offset   int
	filters  map[string]string
}

type updateUserReq struct {
	Passport   *string `json:"passport"`
	Surname    *string `json:"surname"`
	Name       *string `json:"name"`
	Patronymic *string `json:"patronymic"`
	Address    *string `json:"address"`
}

type Repository interface {
	createUser(ctx context.Context, user *user) error
	getUsers(ctx context.Context, searchInfo *searchReq) (*[]user, error)
	removeUser(ctx context.Context, id *userId) error
	updateUser(ctx context.Context, id *userId, updateReq *map[string]interface{}) error
}

type Service interface {
	createUser(ctx context.Context, user *user) error
	getUsers(ctx context.Context, searchInfo *searchReq) (*[]user, error)
	removeUser(ctx context.Context, id *userId) error
	updateUser(ctx context.Context, id *userId, user *updateUserReq) error
}
