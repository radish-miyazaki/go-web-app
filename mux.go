package main

import (
	"context"
	"net/http"

	"github.com/radish-miyazaki/go-web-app/auth"

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

	// バリデーション
	v := validator.New()

	// DBインスタンス
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	// KVSインスタンス
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	// JWTインスタンス
	clocker := clock.RealClocker{}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// Repositoryインスタンス
	r := store.Repository{Clocker: clocker}

	at := &handler.AddTask{Validator: v, Service: &service.AddTask{
		DB:   db,
		Repo: &r,
	}}

	lt := &handler.ListTask{Service: &service.ListTask{
		DB:   db,
		Repo: &r,
	}}

	ru := &handler.RegisterUser{Validator: v, Service: &service.RegisterUser{
		DB:   db,
		Repo: &r,
	}}

	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           &r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	mux.Post("/login", l.ServeHTTP)

	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", at.ServeHTTP)
		r.Get("/", lt.ServeHTTP)
	})

	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application.json charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})

	return mux, cleanup, nil
}
