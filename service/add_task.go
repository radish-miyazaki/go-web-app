package service

import (
	"context"
	"fmt"

	"github.com/radish-miyazaki/go-web-app/auth"

	"github.com/radish-miyazaki/go-web-app/entity"
	"github.com/radish-miyazaki/go-web-app/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	t := &entity.Task{
		UserID: id,
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %v", err)
	}
	return t, nil
}
