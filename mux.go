package main

import (
	"context"
	"net/http"

	"github.com/radish-miyazaki/go-web-app/config"
	"github.com/radish-miyazaki/go-web-app/service"
	"github.com/radish-miyazaki/go-web-app/store/clock"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/radish-miyazaki/go-web-app/handler"
	"github.com/radish-miyazaki/go-web-app/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status": "OK"}`))
	})
	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	r := store.Repository{Clocker: clock.RealClocker{}}
	at := &handler.AddTask{Validator: v, Service: &service.AddTask{
		DB:   db,
		Repo: &r,
	}}
	mux.Post("/tasks", at.ServeHTTP)

	lt := &handler.ListTask{Service: &service.ListTask{
		DB:   db,
		Repo: &r,
	}}
	mux.Get("/tasks", lt.ServeHTTP)

	ru := &handler.RegisterUser{Validator: v, Service: &service.RegisterUser{
		DB:   db,
		Repo: &r,
	}}
	mux.Post("/register", ru.ServeHTTP)

	return mux, cleanup, nil
}
