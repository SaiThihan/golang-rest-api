package routes

import (
	"github.com/SaiThihan/go-basic/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.With(app.Middleware.RequireAuthenticatedUser).Get("/posts/{id}", app.PostHandler.HandleGetPostById)
		r.With(app.Middleware.RequireAuthenticatedUser).Get("/posts", app.PostHandler.HandleGetPosts)
		r.With(app.Middleware.RequireAuthenticatedUser).Post("/posts", app.PostHandler.HandleCreatePost)
		r.With(app.Middleware.RequireAuthenticatedUser).Put("/posts/{id}", app.PostHandler.HandleUpdatePost)
		r.With(app.Middleware.RequireAuthenticatedUser).Delete("/posts/{id}", app.PostHandler.HandleDeletePost)
	})

	r.Get("/health", app.HealthCheck)

	r.Post("/register", app.UserHandler.HandleRegister)
	r.Post("/auth/token", app.TokenHandler.HandleTokenCreate)
	return r
}
