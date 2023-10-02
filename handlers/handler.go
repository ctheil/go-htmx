package handlers

import (
	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
)

func Handler(r *chi.Mux) {

	r.Use(chimiddle.StripSlashes)

	r.Route("/todo", func(router chi.Router) {
		// router.Use(middleware.Authorization)

		router.Get("/", GetTodos)
	})
}
