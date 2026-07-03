package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/SaiThihan/go-basic/internal/app"
	"github.com/SaiThihan/go-basic/internal/routes"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	app.Logger.Println("Go Basic API Server is starting...")

	var port int
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()

	// Setup routes
	r := routes.SetupRoutes(app)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.Logger.Printf("Server is listening on port %d", port)

	err = s.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}
