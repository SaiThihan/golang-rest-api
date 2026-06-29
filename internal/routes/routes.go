package routes

import (
	"github.com/SaiThihan/go-basic/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Get("/posts/{id}", app.PostHandler.HandleGetPostById)
	r.Get("/posts", app.PostHandler.HandleGetPosts)
	r.Post("/posts", app.PostHandler.HandleCreatePost)
	r.Put("/posts/{id}", app.PostHandler.HandleUpdatePost)
	r.Delete("/posts/{id}", app.PostHandler.HandleDeletePost)
	return r
}
