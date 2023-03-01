package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/radish-miyazaki/go-web-app/testutil"
)

func TestLogin_ServeHTTP(t *testing.T) {
	type moq struct {
		token string
		err   error
	}
	type want struct {
		status   int
		respFile string
	}
	tests := map[string]struct {
		reqFile string
		moq     moq
		want    want
	}{
		"ok": {
			reqFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				token: "from_moq",
			},
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/login/ok_resp.json.golden",
			},
		},
		"bad request": {
			reqFile: "testdata/login/bad_request_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/login/bad_request_resp.json.golden",
			},
		},
		"internal server error": {
			reqFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				err: errors.New("error from mock"),
			},
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/login/internal_server_error_resp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))

			moq := &LoginServiceMock{}
			moq.LoginFunc = func(ctx context.Context, name string, password string) (string, error) {
				return tt.moq.token, tt.moq.err
			}
			sut := Login{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.respFile))
		})
	}
}
