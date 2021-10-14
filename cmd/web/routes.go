package main

import (
	"net/http"

	"github.com/bartvanbenthem/gofound-webapp/internal/config"
	"github.com/bartvanbenthem/gofound-webapp/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(NoSurf)
	router.Use(SessionLoad)

	router.Get("/", handlers.Repo.Home)
	router.Get("/home", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)
	router.Get("/items", handlers.Repo.Items)
	router.Get("/contact", handlers.Repo.Contact)
	router.Get("/testform", handlers.Repo.TestForm)
	router.Post("/testform", handlers.Repo.PostTestForm)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}
