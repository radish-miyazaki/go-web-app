package service

import (
	"context"
	"fmt"

	"github.com/radish-miyazaki/go-web-app/entity"
	"github.com/radish-miyazaki/go-web-app/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	ts, err := l.Repo.ListTasks(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %v", err)
	}
	return ts, nil
}
