package service

import (
	"context"
	"fmt"

	"github.com/radish-miyazaki/go-web-app/auth"

	"github.com/radish-miyazaki/go-web-app/entity"
	"github.com/radish-miyazaki/go-web-app/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	ts, err := l.Repo.ListTasks(ctx, l.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %v", err)
	}
	return ts, nil
}
