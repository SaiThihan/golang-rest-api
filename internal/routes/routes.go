package routes

import (
	"github.com/SaiThihan/go-basic/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Get("/posts/{id}", app.PostHandler.HandleGetPostById)
	r.Post("/posts", app.PostHandler.HandleCreatePost)
	return r
}
