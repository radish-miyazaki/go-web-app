package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/radish-miyazaki/go-web-app/entity"
	"github.com/radish-miyazaki/go-web-app/store"
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
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))
			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)
			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.respFile))
		})
	}
}
