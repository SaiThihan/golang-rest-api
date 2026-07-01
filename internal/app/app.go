package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SaiThihan/go-basic/internal/api"
	"github.com/SaiThihan/go-basic/internal/store"
	"github.com/SaiThihan/go-basic/migrations"
)

type Application struct {
	Logger       *log.Logger
	PostHandler  *api.PostHandler
	UserHandler  *api.UserHandler
	TokenHandler *api.TokenHandler
	DB           *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	postgresDB, err := store.Open()

	if err != nil {
		return nil, err
	}

	err = store.MigrateFs(postgresDB, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	postStore := store.NewPostgresPostStore(postgresDB)
	postHandler := api.NewPostHandler(postStore, logger)
	userStore := store.NewPostgresUserStore(postgresDB)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenStore := store.NewPostgresTokenStore(postgresDB)
	tokenHandler := api.NewTokenHandler(userStore, tokenStore, logger)

	app := &Application{
		Logger:       logger,
		PostHandler:  postHandler,
		UserHandler:  userHandler,
		TokenHandler: tokenHandler,
		DB:           postgresDB,
	}
	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is Ok\n")
}
