package store

import (
	"context"
	"testing"

	"github.com/radish-miyazaki/go-web-app/testutil/fixture"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/radish-miyazaki/go-web-app/entity"
	"github.com/radish-miyazaki/go-web-app/store/clock"
	"github.com/radish-miyazaki/go-web-app/testutil"
)

// 実際のDBを使ったテストコード
func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() {
		_ = tx.Rollback()
	})
	if err != nil {
		t.Fatal(err)
	}
	wantUserID, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, wantUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		UserID:    entity.UserID(1),
		Title:     "ok task",
		Status:    "todo",
		CreatedAt: c.Now(),
		UpdatedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	mock.ExpectExec(`INSERT INTO task \(user_id, title, status, created_at, updated_at\) VALUES \(\?, \?, \?, \?, \?\)`).
		WithArgs(okTask.UserID, okTask.Title, okTask.Status, okTask.CreatedAt, okTask.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()

	u := fixture.User(nil)
	result, err := db.ExecContext(ctx,
		"INSERT INTO user (name, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		u.Name,
		u.Password,
		u.Role,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		t.Fatalf("insert user %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("got user_id: %v", err)
	}
	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, db Execer) (entity.UserID, entity.Tasks) {
	t.Helper()
	userID := prepareUser(ctx, t, db)
	otherUserID := prepareUser(ctx, t, db)

	if _, err := db.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title:     "want task 1",
			Status:    "todo",
			CreatedAt: c.Now(),
			UpdatedAt: c.Now(),
			UserID:    userID,
		},
		{
			Title:     "want task 2",
			Status:    "done",
			CreatedAt: c.Now(),
			UpdatedAt: c.Now(),
			UserID:    userID,
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			Title:     "not want",
			Status:    "todo",
			CreatedAt: c.Now(),
			UpdatedAt: c.Now(),
			UserID:    otherUserID,
		},
		wants[1],
	}

	result, err := db.ExecContext(
		ctx,
		"INSERT INTO task (user_id, title, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?), (?, ?, ?, ?, ?);",
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].CreatedAt, tasks[0].UpdatedAt,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].CreatedAt, tasks[1].UpdatedAt,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].CreatedAt, tasks[2].UpdatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)
	return userID, wants
}
