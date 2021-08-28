package main

import (
	"net/http"

	"github.com/bartvanbenthem/gofound-web/internal/config"
	"github.com/bartvanbenthem/gofound-web/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/home", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/items", handlers.Repo.Items)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Post("/testpost", handlers.Repo.PostTest)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
