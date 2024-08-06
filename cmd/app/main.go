package main

import (
	"fmt"
	"net/http"
	"os"
	"tasks/internal/database"
	"tasks/internal/transport"
	"tasks/pkg/auth"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Info().Msg("Initializing database")
}

func main() {
	connString := os.Getenv("POSTGRES_CONN_STRING")
	fmt.Print(connString)

	storage, err := database.NewPostgresStorage(connString)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}

	mux := http.NewServeMux()

	auth := auth.Auth{
		S: storage,
	}

	res := transport.NewResource(storage)

	mux.HandleFunc("POST /users", res.CreateUser)
	mux.HandleFunc("GET /tasks", auth.CheckAuth(res.GetTasks))
	mux.HandleFunc("POST /tasks", auth.CheckAuth(res.CreateTask))
	mux.HandleFunc("PUT /tasks/{id}", auth.CheckAuth(res.UpdateTask))
	mux.HandleFunc("DELETE /tasks/{id}", auth.CheckAuth(res.DeleteTask))

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Error().Err(err).Msg("Failed to listen and serve")
	}
}
