package store

import (
	"context"

	"github.com/radish-miyazaki/go-web-app/entity"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer, id entity.UserID) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := "SELECT id, user_id, title, status, created_at, updated_at FROM task WHERE user_id = ?;"
	if err := db.SelectContext(ctx, &tasks, sql, id); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.CreatedAt = r.Clocker.Now()
	t.UpdatedAt = r.Clocker.Now()
	sql := "INSERT INTO task (user_id, title, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	result, err := db.ExecContext(ctx, sql, t.UserID, t.Title, t.Status, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}
