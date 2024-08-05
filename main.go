package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Info().Msg("Initializing database")
}

func main() {
	connString := os.Getenv("POSTGRES_CONN_STRING")
	fmt.Print(connString)

	storage, err := NewStorage(connString)

	if err != nil {
		log.Error().Err(err).Msg("Failed to open database")
	}

	mux := http.NewServeMux()

	tasks := TasksResource{
		Storage: storage,
	}

	auth := Auth{
		Storage: storage,
	}

	users := UsersResource{
		Storage: storage,
	}

	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("GET /tasks", auth.checkAuth(tasks.GetTasks))
	mux.HandleFunc("POST /tasks", auth.checkAuth(tasks.CreateTask))
	mux.HandleFunc("PUT /tasks/{id}", auth.checkAuth(tasks.UpdateTask))
	mux.HandleFunc("DELETE /tasks/{id}", auth.checkAuth(tasks.DeleteTask))

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Error().Err(err).Msg("Failed to listen and serve")
	}
}

type TasksResource struct {
	Storage *Storage
}

func (tr *TasksResource) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := tr.Storage.GetTasks()
	fmt.Print(tasks)

	if err != nil {
		log.Error().Err(err).Msg("Error getting tasks")
	}

	err = json.NewEncoder(w).Encode(tasks)

	if err != nil {
		fmt.Println("Failed to encode: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (tr *TasksResource) CreateTask(w http.ResponseWriter, r *http.Request) {
	var reqBody Task

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		log.Error().Err(err).Msg("Failed to decode")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = tr.Storage.InsertTask(reqBody)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (tr *TasksResource) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var reqBody Task

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		fmt.Println("Failed to encode: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = tr.Storage.UpdateTask(id, reqBody)

	if err != nil {
		log.Error().Err(err).Msg("Failed to update task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (tr *TasksResource) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := tr.Storage.DeleteTask(id)

	if err != nil {
		log.Error().Err(err).Msg("Error deleting task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type UsersResource struct {
	Storage *Storage
}

func (ur *UsersResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Error().Err(err).Msg("Failed to decode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ur.Storage.InsertUser(user)

	if err != nil {
		log.Error().Err(err).Msg("Error creating user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
