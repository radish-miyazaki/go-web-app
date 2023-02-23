package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/radish-miyazaki/go-web-app/entity"
	"github.com/radish-miyazaki/go-web-app/testutil"
)

func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status   int
		respFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/add_task/ok_resp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_request_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/add_task/bad_request_resp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))
			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)
			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.respFile))
		})
	}
}
