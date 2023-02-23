package service

import (
	"context"

	"github.com/radish-miyazaki/go-web-app/entity"
	"github.com/radish-miyazaki/go-web-app/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskListener UserRegister
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}
